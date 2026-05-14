// SPDX-License-Identifier: GPL-3.0-only

package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/cinmou/ClipBridgeServer/internal/api"
	"github.com/cinmou/ClipBridgeServer/internal/cleanup"
	"github.com/cinmou/ClipBridgeServer/internal/config"
	"github.com/cinmou/ClipBridgeServer/internal/store"
	"github.com/cinmou/ClipBridgeServer/internal/webdav"
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

	if err := prepareConfig(cfg, log.Default()); err != nil {
		log.Fatalf("prepare configuration failed: %v", err)
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

	webdavService := webdav.NewService(dbStore, cfg)

	router := api.NewRouter(dbStore, cfg, cleanupService, webdavService, webui.Handler())
	addr := cfg.Server.Address()

	log.Printf("ClipBridgeServer starting on %s", addr)
	log.Printf("ClipBridgeServer sqlite ready at %s", cfg.Storage.DatabasePath)
	logStartup(log.Default(), cfg)

	if err := serve(cfg, router); err != nil {
		logSafePrintf(log.Default(), "server stopped: %v", err)
		os.Exit(1)
	}
}
