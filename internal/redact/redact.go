// SPDX-License-Identifier: GPL-3.0-only

package redact

import "regexp"

var sensitivePatterns = []struct {
	pattern *regexp.Regexp
	replace string
}{
	{regexp.MustCompile(`cb_admin_[A-Za-z0-9_-]+`), `cb_admin_[REDACTED]`},
	{regexp.MustCompile(`cb_device_[A-Za-z0-9_-]+`), `cb_device_[REDACTED]`},
	{regexp.MustCompile(`(?i)(Authorization:\s*Bearer\s+)([^\s]+)`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)(Cookie:\s*)(.+)`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)(pairing_code=)([^&\s]+)`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)("pairing_code"\s*:\s*")([^"]+)(")`), `${1}[REDACTED]${3}`},
	{regexp.MustCompile(`(?i)(pairing_code\s*[:=]\s*)([^\s,}]+)`), `${1}[REDACTED]`},
	{regexp.MustCompile(`(?i)("password"\s*:\s*")([^"]+)(")`), `${1}[REDACTED]${3}`},
	{regexp.MustCompile(`(?i)(password\s*[:=]\s*)([^\s,}]+)`), `${1}[REDACTED]`},
}

// Text redacts common secret-bearing values before log emission.
func Text(raw string) string {
	redacted := raw
	for _, rule := range sensitivePatterns {
		redacted = rule.pattern.ReplaceAllString(redacted, rule.replace)
	}
	return redacted
}
