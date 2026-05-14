// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cinmou/ClipBridgeServer/internal/auth"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/redact"
)

const insecureBindWarning = "Warning: ClipBridgeServer is listening on 0.0.0.0 without TLS. Tokens may be exposed on untrusted networks."

var (
	listenAndServe    = http.ListenAndServe
	listenAndServeTLS = http.ListenAndServeTLS
)

func prepareConfig(cfg *config.Config, logger *log.Logger) error {
	token, created, err := resolveAdminToken(cfg.Storage.DataDir, cfg.Auth.Token)
	if err != nil {
		return err
	}

	cfg.Auth.Token = token
	if created && logger != nil {
		logger.Printf("ClipBridgeServer generated admin token: %s", token)
	}

	return nil
}

func resolveAdminToken(dataDir, configuredToken string) (string, bool, error) {
	if token := strings.TrimSpace(configuredToken); token != "" {
		return token, false, nil
	}

	secretsDir := filepath.Join(dataDir, "secrets")
	if err := os.MkdirAll(secretsDir, 0o700); err != nil {
		return "", false, fmt.Errorf("create secrets directory %q: %w", secretsDir, err)
	}

	tokenPath := filepath.Join(secretsDir, "admin_token")
	if data, err := os.ReadFile(tokenPath); err == nil {
		if token := strings.TrimSpace(string(data)); token != "" {
			return token, false, nil
		}
	} else if !os.IsNotExist(err) {
		return "", false, fmt.Errorf("read admin token file %q: %w", tokenPath, err)
	}

	token, err := auth.GenerateAdminToken()
	if err != nil {
		return "", false, err
	}

	if err := os.WriteFile(tokenPath, []byte(token+"\n"), 0o600); err != nil {
		return "", false, fmt.Errorf("write admin token file %q: %w", tokenPath, err)
	}

	return token, true, nil
}

func startupWarning(cfg *config.Config) string {
	if cfg.Server.Host == "0.0.0.0" && !cfg.TLS.Enabled {
		return insecureBindWarning
	}
	return ""
}

func logStartup(logger *log.Logger, cfg *config.Config) {
	if logger == nil {
		return
	}

	if warning := startupWarning(cfg); warning != "" {
		logger.Print(warning)
	}
}

func serve(cfg *config.Config, handler http.Handler) error {
	addr := cfg.Server.Address()
	if cfg.TLS.Enabled {
		return listenAndServeTLS(addr, cfg.TLS.CertFile, cfg.TLS.KeyFile, handler)
	}
	return listenAndServe(addr, handler)
}

func logSafePrintf(logger *log.Logger, format string, args ...any) {
	if logger == nil {
		return
	}
	logger.Print(redact.Text(fmt.Sprintf(format, args...)))
}
