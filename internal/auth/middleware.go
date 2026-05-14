// SPDX-License-Identifier: GPL-3.0-only

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const pairingCodeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// ExtractBearerToken reads the Authorization header and returns the bearer token
// value when present. Keeping this parsing logic centralized avoids subtle
// differences between admin and device authentication paths.
func ExtractBearerToken(req *http.Request) (string, bool) {
	headerValue := req.Header.Get("Authorization")
	if !strings.HasPrefix(headerValue, "Bearer ") {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(headerValue, "Bearer "))
	if token == "" {
		return "", false
	}

	return token, true
}

// HashSecret stores opaque credentials as deterministic SHA-256 hashes instead
// of plaintext values. This keeps pairing codes and device tokens out of the
// database while still allowing equality checks during authentication.
func HashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// GeneratePairingCode creates a short human-friendly pairing code for manual
// device onboarding.
func GeneratePairingCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("pairing code length must be greater than 0")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate pairing code bytes: %w", err)
	}

	builder := strings.Builder{}
	builder.Grow(length)
	for _, value := range bytes {
		builder.WriteByte(pairingCodeAlphabet[int(value)%len(pairingCodeAlphabet)])
	}

	return builder.String(), nil
}

// GenerateDeviceToken creates a long opaque token that clients can store and
// reuse for future authenticated requests.
func GenerateDeviceToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate device token bytes: %w", err)
	}

	return "cb_device_" + base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateAdminToken creates a long opaque admin token for management APIs.
func GenerateAdminToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate admin token bytes: %w", err)
	}

	return "cb_admin_" + base64.RawURLEncoding.EncodeToString(bytes), nil
}
