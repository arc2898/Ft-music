// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/arc2898/ft-music/go/internal/repository/sqlite"
	"github.com/arc2898/ft-music/go/internal/server"
	"github.com/arc2898/ft-music/go/internal/service"
	"github.com/arc2898/ft-music/go/internal/streaming"
)

func main() {
	addr := envOrDefault("FT_MUSIC_ADDR", "127.0.0.1:8080")
	apiToken := strings.TrimSpace(os.Getenv("FT_MUSIC_API_TOKEN"))
	if isNonLoopbackAddress(addr) && apiToken == "" {
		log.Fatal("FT_MUSIC_API_TOKEN is required when FT_MUSIC_ADDR exposes the API beyond loopback")
	}
	dbPath := envOrDefault("FT_MUSIC_DB", "ft-music.db")

	tracks, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatalf("FT-music could not start: %v", err)
	}
	defer func() {
		if err := tracks.Close(); err != nil {
			log.Printf("FT-music database close failed: %v", err)
		}
	}()

	music := service.NewMusic(tracks, streaming.NewCatalog())
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.New(music, server.WithAPIKey(apiToken)).Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("FT-music API listening on %s (database: %s)", addr, dbPath)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("FT-music API stopped: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		log.Printf("FT-music API is shutting down")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("FT-music graceful shutdown failed: %v", err)
			if closeErr := httpServer.Close(); closeErr != nil {
				log.Printf("FT-music forced shutdown failed: %v", closeErr)
			}
		}
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func isNonLoopbackAddress(address string) bool {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return true
	}
	if host == "localhost" {
		return false
	}
	ip := net.ParseIP(host)
	return ip == nil || !ip.IsLoopback()
}
