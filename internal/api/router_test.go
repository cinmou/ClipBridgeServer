// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestLoggingMiddlewareDoesNotLeakSensitiveHeadersOrQueryValues(t *testing.T) {
	testContext := newTestContext(t)

	var logs bytes.Buffer
	originalWriter := log.Writer()
	originalFlags := log.Flags()
	log.SetOutput(&logs)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(originalWriter)
		log.SetFlags(originalFlags)
	})

	request := performRequest(
		t,
		testContext.router,
		http.MethodGet,
		"/api/health?pairing_code=ABCDEFGH&password=super-secret",
		nil,
		"cb_admin_super_secret",
	)
	if request.Code != http.StatusOK {
		t.Fatalf("GET /api/health status = %d, want %d", request.Code, http.StatusOK)
	}

	logOutput := logs.String()
	for _, forbidden := range []string{
		"cb_admin_super_secret",
		"sessionid=abc123",
		"ABCDEFGH",
		"super-secret",
		"Authorization:",
		"Cookie:",
	} {
		if strings.Contains(logOutput, forbidden) {
			t.Fatalf("log output leaked %q: %q", forbidden, logOutput)
		}
	}
	if !strings.Contains(logOutput, "pairing_code=[REDACTED]") {
		t.Fatalf("log output = %q, want redacted pairing code", logOutput)
	}
}
