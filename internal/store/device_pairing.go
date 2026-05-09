// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/auth"
)

var (
	// ErrAlreadyUsed marks one-time pairing codes that have already been
	// consumed.
	ErrAlreadyUsed = errors.New("store: pairing code already used")
	// ErrExpired marks pairing codes that are no longer valid because their
	// expiry timestamp has passed.
	ErrExpired = errors.New("store: pairing code expired")
)

// Device describes one paired client device. The token itself never appears
// here; only metadata that is safe to return through the API.
type Device struct {
	ID         int64
	Name       string
	CreatedAt  string
	LastSeenAt string
	RevokedAt  string
}

// PairingCode describes a short-lived one-time code shown to the user during
// device onboarding.
type PairingCode struct {
	Code      string
	ExpiresAt string
}

// CreatePairingCode inserts a one-time pairing code that expires at the given
// time. The database only stores the hash of the code, never the plaintext.
func (s *SQLiteStore) CreatePairingCode(ctx context.Context, expiresAt time.Time) (*PairingCode, error) {
	for range 5 {
		rawCode, err := auth.GeneratePairingCode(8)
		if err != nil {
			return nil, err
		}

		expiresAtText := expiresAt.UTC().Format(time.RFC3339)
		_, err = s.db.ExecContext(
			ctx,
			`INSERT INTO pairing_codes (code_hash, expires_at) VALUES (?, ?)`,
			auth.HashSecret(rawCode),
			expiresAtText,
		)
		if err != nil {
			if isUniqueConstraint(err) {
				continue
			}
			return nil, fmt.Errorf("insert pairing code: %w", err)
		}

		return &PairingCode{
			Code:      rawCode,
			ExpiresAt: expiresAtText,
		}, nil
	}

	return nil, fmt.Errorf("create pairing code: too many collisions")
}

// ExchangePairingCode validates and consumes one pairing code, then creates a
// device token for the new client. The code becomes unusable immediately after
// a successful exchange.
func (s *SQLiteStore) ExchangePairingCode(ctx context.Context, rawCode string, deviceName string) (*Device, string, error) {
	codeHash := auth.HashSecret(strings.TrimSpace(rawCode))
	now := time.Now().UTC()
	nowText := now.Format(time.RFC3339)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("begin pairing exchange transaction: %w", err)
	}

	var pairingCodeID int64
	var expiresAtText string
	var usedAt sql.NullString

	err = tx.QueryRowContext(
		ctx,
		`SELECT id, expires_at, used_at
		 FROM pairing_codes
		 WHERE code_hash = ?
		 LIMIT 1`,
		codeHash,
	).Scan(&pairingCodeID, &expiresAtText, &usedAt)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("query pairing code: %w", err)
	}

	if usedAt.Valid {
		_ = tx.Rollback()
		return nil, "", ErrAlreadyUsed
	}

	expiresAt, err := time.Parse(time.RFC3339, expiresAtText)
	if err != nil {
		_ = tx.Rollback()
		return nil, "", fmt.Errorf("parse pairing code expiry: %w", err)
	}
	if !now.Before(expiresAt) {
		_ = tx.Rollback()
		return nil, "", ErrExpired
	}

	rawDeviceToken := ""
	tokenHash := ""
	inserted := false
	for range 5 {
		rawDeviceToken, err = auth.GenerateDeviceToken()
		if err != nil {
			_ = tx.Rollback()
			return nil, "", err
		}

		tokenHash = auth.HashSecret(rawDeviceToken)
		result, execErr := tx.ExecContext(
			ctx,
			`INSERT INTO devices (name, token_hash, created_at, last_seen_at)
			 VALUES (?, ?, ?, ?)`,
			normalizeDeviceName(deviceName),
			tokenHash,
			nowText,
			nowText,
		)
		if execErr != nil {
			if isUniqueConstraint(execErr) {
				continue
			}
			_ = tx.Rollback()
			return nil, "", fmt.Errorf("insert paired device: %w", execErr)
		}

		deviceID, execErr := result.LastInsertId()
		if execErr != nil {
			_ = tx.Rollback()
			return nil, "", fmt.Errorf("read paired device id: %w", execErr)
		}

		if _, execErr := tx.ExecContext(
			ctx,
			`UPDATE pairing_codes
			 SET used_at = ?, used_by_device_id = ?
			 WHERE id = ?`,
			nowText,
			deviceID,
			pairingCodeID,
		); execErr != nil {
			_ = tx.Rollback()
			return nil, "", fmt.Errorf("mark pairing code as used: %w", execErr)
		}

		if err := tx.Commit(); err != nil {
			return nil, "", fmt.Errorf("commit pairing exchange: %w", err)
		}

		device, err := s.GetDeviceByID(ctx, deviceID)
		if err != nil {
			return nil, "", err
		}

		inserted = true
		return device, rawDeviceToken, nil
	}

	if !inserted {
		_ = tx.Rollback()
		return nil, "", fmt.Errorf("create paired device: too many token collisions")
	}

	return nil, "", fmt.Errorf("create paired device failed")
}

// AuthenticateDeviceToken validates an active device token and updates the
// device last-seen timestamp when the token is accepted.
func (s *SQLiteStore) AuthenticateDeviceToken(ctx context.Context, rawToken string) (*Device, error) {
	nowText := time.Now().UTC().Format(time.RFC3339)
	tokenHash := auth.HashSecret(strings.TrimSpace(rawToken))

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE devices
		 SET last_seen_at = ?
		 WHERE token_hash = ? AND revoked_at IS NULL`,
		nowText,
		tokenHash,
	)
	if err != nil {
		return nil, fmt.Errorf("update device last_seen_at: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read authenticated device row count: %w", err)
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return s.getOneDevice(
		ctx,
		`SELECT id, name, created_at, COALESCE(last_seen_at, ''), COALESCE(revoked_at, '')
		 FROM devices
		 WHERE token_hash = ? AND revoked_at IS NULL
		 LIMIT 1`,
		tokenHash,
	)
}

// ListDevices returns all paired devices, including revoked ones, so the future
// Web UI can show a full audit trail.
func (s *SQLiteStore) ListDevices(ctx context.Context) ([]Device, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name, created_at, COALESCE(last_seen_at, ''), COALESCE(revoked_at, '')
		 FROM devices
		 ORDER BY id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()

	devices := make([]Device, 0)
	for rows.Next() {
		var device Device
		if err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.CreatedAt,
			&device.LastSeenAt,
			&device.RevokedAt,
		); err != nil {
			return nil, fmt.Errorf("scan device row: %w", err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate devices: %w", err)
	}

	return devices, nil
}

// RevokeDevice invalidates one device token without deleting the device record.
func (s *SQLiteStore) RevokeDevice(ctx context.Context, id int64) error {
	nowText := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.ExecContext(
		ctx,
		`UPDATE devices
		 SET revoked_at = ?
		 WHERE id = ? AND revoked_at IS NULL`,
		nowText,
		id,
	)
	if err != nil {
		return fmt.Errorf("revoke device %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read revoked device row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetDeviceByID returns one device by id.
func (s *SQLiteStore) GetDeviceByID(ctx context.Context, id int64) (*Device, error) {
	return s.getOneDevice(
		ctx,
		`SELECT id, name, created_at, COALESCE(last_seen_at, ''), COALESCE(revoked_at, '')
		 FROM devices
		 WHERE id = ?
		 LIMIT 1`,
		id,
	)
}

func (s *SQLiteStore) getOneDevice(ctx context.Context, query string, args ...any) (*Device, error) {
	var device Device

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&device.ID,
		&device.Name,
		&device.CreatedAt,
		&device.LastSeenAt,
		&device.RevokedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query device: %w", err)
	}

	return &device, nil
}

func normalizeDeviceName(deviceName string) string {
	name := strings.TrimSpace(deviceName)
	if name == "" {
		return "Unnamed Device"
	}

	return name
}

func isUniqueConstraint(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique")
}
