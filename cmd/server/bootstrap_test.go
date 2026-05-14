// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cinmou/ClipBridgeServer/internal/config"
)

func TestPrepareConfigGeneratesAdminTokenWhenMissing(t *testing.T) {
	cfg := config.Default()
	cfg.Auth.Token = ""
	cfg.Storage.DataDir = t.TempDir()

	var logs bytes.Buffer
	logger := log.New(&logs, "", 0)

	if err := prepareConfig(cfg, logger); err != nil {
		t.Fatalf("prepareConfig() error = %v", err)
	}

	if !strings.HasPrefix(cfg.Auth.Token, "cb_admin_") {
		t.Fatalf("generated admin token = %q, want cb_admin_ prefix", cfg.Auth.Token)
	}

	weakDefaults := map[string]bool{
		"":          true,
		"change-me": true,
		"admin":     true,
		"password":  true,
	}
	if weakDefaults[cfg.Auth.Token] {
		t.Fatalf("generated admin token = %q, should not be a weak default", cfg.Auth.Token)
	}

	tokenPath := filepath.Join(cfg.Storage.DataDir, "secrets", "admin_token")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", tokenPath, err)
	}
	if strings.TrimSpace(string(data)) != cfg.Auth.Token {
		t.Fatalf("stored admin token = %q, want %q", strings.TrimSpace(string(data)), cfg.Auth.Token)
	}
	if !strings.Contains(logs.String(), cfg.Auth.Token) {
		t.Fatalf("startup log = %q, want generated token printed once", logs.String())
	}

	if runtimeMode := os.FileMode(0o600); runtimeMode != 0 {
		info, err := os.Stat(tokenPath)
		if err != nil {
			t.Fatalf("Stat(%q) error = %v", tokenPath, err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("admin token file mode = %#o, want %#o", got, 0o600)
		}
	}
}

func TestPrepareConfigLoadsExistingAdminTokenWithoutReprinting(t *testing.T) {
	cfg := config.Default()
	cfg.Auth.Token = ""
	cfg.Storage.DataDir = t.TempDir()

	var firstLogs bytes.Buffer
	if err := prepareConfig(cfg, log.New(&firstLogs, "", 0)); err != nil {
		t.Fatalf("first prepareConfig() error = %v", err)
	}
	firstToken := cfg.Auth.Token

	cfg.Auth.Token = ""
	var secondLogs bytes.Buffer
	if err := prepareConfig(cfg, log.New(&secondLogs, "", 0)); err != nil {
		t.Fatalf("second prepareConfig() error = %v", err)
	}

	if cfg.Auth.Token != firstToken {
		t.Fatalf("reloaded admin token = %q, want %q", cfg.Auth.Token, firstToken)
	}
	if secondLogs.Len() != 0 {
		t.Fatalf("second startup logs = %q, want token to not be printed again", secondLogs.String())
	}
}

func TestStartupWarning(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Host = "0.0.0.0"
	cfg.TLS.Enabled = false

	var logs bytes.Buffer
	logStartup(log.New(&logs, "", 0), cfg)

	if !strings.Contains(logs.String(), insecureBindWarning) {
		t.Fatalf("startup log = %q, want insecure bind warning", logs.String())
	}
}

func TestServeUsesHTTPWhenTLSDisabled(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8787
	cfg.TLS.Enabled = false

	originalHTTP := listenAndServe
	originalTLS := listenAndServeTLS
	defer func() {
		listenAndServe = originalHTTP
		listenAndServeTLS = originalTLS
	}()

	calledHTTP := false
	calledTLS := false
	listenAndServe = func(addr string, handler http.Handler) error {
		calledHTTP = true
		if addr != "127.0.0.1:8787" {
			t.Fatalf("ListenAndServe addr = %q, want %q", addr, "127.0.0.1:8787")
		}
		return errors.New("stop")
	}
	listenAndServeTLS = func(string, string, string, http.Handler) error {
		calledTLS = true
		return nil
	}

	if err := serve(cfg, http.NotFoundHandler()); err == nil || err.Error() != "stop" {
		t.Fatalf("serve() error = %v, want stop", err)
	}
	if !calledHTTP || calledTLS {
		t.Fatalf("serve() dispatch http=%v tls=%v, want http=true tls=false", calledHTTP, calledTLS)
	}
}

func TestServeUsesTLSWhenEnabled(t *testing.T) {
	cfg := config.Default()
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8787
	cfg.TLS.Enabled = true
	cfg.TLS.CertFile = "cert.pem"
	cfg.TLS.KeyFile = "key.pem"

	originalHTTP := listenAndServe
	originalTLS := listenAndServeTLS
	defer func() {
		listenAndServe = originalHTTP
		listenAndServeTLS = originalTLS
	}()

	calledHTTP := false
	calledTLS := false
	listenAndServe = func(string, http.Handler) error {
		calledHTTP = true
		return nil
	}
	listenAndServeTLS = func(addr, certFile, keyFile string, handler http.Handler) error {
		calledTLS = true
		if addr != "127.0.0.1:8787" || certFile != "cert.pem" || keyFile != "key.pem" {
			t.Fatalf("ListenAndServeTLS args = (%q, %q, %q), want (%q, %q, %q)", addr, certFile, keyFile, "127.0.0.1:8787", "cert.pem", "key.pem")
		}
		return errors.New("stop")
	}

	if err := serve(cfg, http.NotFoundHandler()); err == nil || err.Error() != "stop" {
		t.Fatalf("serve() error = %v, want stop", err)
	}
	if calledHTTP || !calledTLS {
		t.Fatalf("serve() dispatch http=%v tls=%v, want http=false tls=true", calledHTTP, calledTLS)
	}
}
