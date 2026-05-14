// SPDX-License-Identifier: GPL-3.0-only

package redact

import (
	"strings"
	"testing"
)

func TestTextRedactsSensitiveValues(t *testing.T) {
	raw := strings.Join([]string{
		"Authorization: Bearer cb_admin_secret_admin_token",
		"Cookie: sessionid=abc123; other=value",
		"cb_admin_secret_admin_token",
		"cb_device_secret_device_token",
		`{"pairing_code":"ABCDEFGH","password":"super-secret"}`,
		"pairing_code=ABCDEFGH",
		"password=super-secret",
	}, "\n")

	redacted := Text(raw)

	for _, secret := range []string{
		"secret_admin_token",
		"secret_device_token",
		"ABCDEFGH",
		"super-secret",
		"sessionid=abc123",
	} {
		if strings.Contains(redacted, secret) {
			t.Fatalf("redacted text still contains %q: %q", secret, redacted)
		}
	}

	for _, want := range []string{
		"Authorization: Bearer [REDACTED]",
		"Cookie: [REDACTED]",
		"cb_admin_[REDACTED]",
		"cb_device_[REDACTED]",
		`"pairing_code":"[REDACTED]"`,
		`"password":"[REDACTED]"`,
		"pairing_code=[REDACTED]",
		"password=[REDACTED]",
	} {
		if !strings.Contains(redacted, want) {
			t.Fatalf("redacted text = %q, want substring %q", redacted, want)
		}
	}
}
