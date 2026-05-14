// SPDX-License-Identifier: GPL-3.0-only

package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cinmou/ClipBridgeServer/internal/auth"
)

func TestPairingCodeLifecycleAndHashStorage(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "clipbridge.db")
	dbStore, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		_ = dbStore.Close()
	})

	ctx := context.Background()

	pairingCode, err := dbStore.CreatePairingCode(ctx, time.Now().UTC().Add(5*time.Minute))
	if err != nil {
		t.Fatalf("CreatePairingCode() error = %v", err)
	}

	device, deviceToken, err := dbStore.ExchangePairingCode(ctx, pairingCode.Code, "MacBook")
	if err != nil {
		t.Fatalf("ExchangePairingCode() error = %v", err)
	}
	if device.Name != "MacBook" {
		t.Fatalf("paired device name = %q, want %q", device.Name, "MacBook")
	}
	if deviceToken == "" {
		t.Fatalf("device token should not be empty")
	}
	if !strings.HasPrefix(deviceToken, "cb_device_") {
		t.Fatalf("device token = %q, want cb_device_ prefix", deviceToken)
	}

	authenticatedDevice, err := dbStore.AuthenticateDeviceToken(ctx, deviceToken)
	if err != nil {
		t.Fatalf("AuthenticateDeviceToken() error = %v", err)
	}
	if authenticatedDevice.ID != device.ID {
		t.Fatalf("authenticated device id = %d, want %d", authenticatedDevice.ID, device.ID)
	}

	_, _, err = dbStore.ExchangePairingCode(ctx, pairingCode.Code, "iPhone")
	if !errors.Is(err, ErrAlreadyUsed) {
		t.Fatalf("reuse pairing code error = %v, want ErrAlreadyUsed", err)
	}

	if err := dbStore.RevokeDevice(ctx, device.ID); err != nil {
		t.Fatalf("RevokeDevice() error = %v", err)
	}

	_, err = dbStore.AuthenticateDeviceToken(ctx, deviceToken)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("AuthenticateDeviceToken() after revoke error = %v, want ErrNotFound", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	var storedCodeHash string
	if err := db.QueryRow("SELECT code_hash FROM pairing_codes LIMIT 1").Scan(&storedCodeHash); err != nil {
		t.Fatalf("query pairing code hash error = %v", err)
	}
	if storedCodeHash != auth.HashSecret(pairingCode.Code) {
		t.Fatalf("stored code hash = %q, want hash of pairing code", storedCodeHash)
	}
	if storedCodeHash == pairingCode.Code {
		t.Fatalf("stored code hash should not equal plaintext pairing code")
	}

	var storedTokenHash string
	if err := db.QueryRow("SELECT token_hash FROM devices WHERE id = ?", device.ID).Scan(&storedTokenHash); err != nil {
		t.Fatalf("query token hash error = %v", err)
	}
	if storedTokenHash != auth.HashSecret(deviceToken) {
		t.Fatalf("stored token hash = %q, want hash of device token", storedTokenHash)
	}
	if storedTokenHash == deviceToken {
		t.Fatalf("stored token hash should not equal plaintext device token")
	}
}

func TestExpiredPairingCodeFails(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "clipbridge.db")
	dbStore, err := OpenSQLite(dbPath)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		_ = dbStore.Close()
	})

	ctx := context.Background()
	pairingCode, err := dbStore.CreatePairingCode(ctx, time.Now().UTC().Add(-1*time.Minute))
	if err != nil {
		t.Fatalf("CreatePairingCode(expired) error = %v", err)
	}

	_, _, err = dbStore.ExchangePairingCode(ctx, pairingCode.Code, "Expired Device")
	if !errors.Is(err, ErrExpired) {
		t.Fatalf("ExchangePairingCode(expired) error = %v, want ErrExpired", err)
	}
}
