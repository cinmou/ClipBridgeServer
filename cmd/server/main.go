// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/cinmou/ClipBridgeServer/internal/api"
	"github.com/cinmou/ClipBridgeServer/internal/cleanup"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
	webui "github.com/cinmou/ClipBridgeServer/web"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("configuration file %q not found", *configPath)
		}
		log.Fatalf("load configuration failed: %v", err)
	}

	dbStore, err := store.OpenSQLite(cfg.Storage.DatabasePath)
	if err != nil {
		log.Fatalf("initialize sqlite failed: %v", err)
	}
	defer func() {
		if err := dbStore.Close(); err != nil {
			log.Printf("close sqlite failed: %v", err)
		}
	}()

	cleanupService, err := cleanup.NewService(dbStore, cfg)
	if err != nil {
		log.Fatalf("initialize cleanup service failed: %v", err)
	}
	defer cleanupService.Close()

	router := api.NewRouter(dbStore, cfg, cleanupService, webui.Handler())
	addr := cfg.Server.Address()

	log.Printf("ClipBridgeServer starting on %s", addr)
	log.Printf("ClipBridgeServer sqlite ready at %s", cfg.Storage.DatabasePath)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
