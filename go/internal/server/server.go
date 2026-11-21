// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package server

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
	"github.com/arc2898/ft-music/go/internal/service"
	"github.com/arc2898/ft-music/go/internal/streaming"
)

const (
	maxJSONBody  = 1 << 20
	defaultLimit = 50
	maxLimit     = 100
	maxQuerySize = 256
)

type Server struct {
	music  *service.Music
	apiKey string
}

type Option func(*Server)

// WithAPIKey protects application endpoints with an Authorization: Bearer token.
func WithAPIKey(key string) Option {
	return func(s *Server) { s.apiKey = strings.TrimSpace(key) }
}

func New(music *service.Music, options ...Option) *Server {
	s := &Server{music: music}
	for _, option := range options {
		if option != nil {
			option(s)
		}
	}
	return s
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("GET /api/v1/health", s.health)
	mux.HandleFunc("GET /api/v1/library", s.listLibrary)
	mux.HandleFunc("GET /api/v1/library/{id}", s.getTrack)
	mux.HandleFunc("POST /api/v1/library", s.addTrack)
	mux.HandleFunc("DELETE /api/v1/library/{id}", s.deleteTrack)
	mux.HandleFunc("GET /api/v1/search", s.search)
	mux.HandleFunc("GET /api/v1/providers", s.listProviders)
	mux.HandleFunc("GET /api/v1/providers/{provider}/tracks/{id}/stream", s.streamURL)

	var handler http.Handler = mux
	handler = recoverMiddleware(handler)
	handler = securityHeaders(handler)
	if s.apiKey != "" {
		handler = s.apiKeyMiddleware(handler)
	}
	handler = requestIDMiddleware(handler)
	return handler
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "ft-music-api",
		"brand":   "FT-music",
		"status":  "ok",
	})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if s.music == nil {
		writeError(w, http.StatusServiceUnavailable, "not_ready", "music service is not configured")
		return
	}
	if err := s.music.Ready(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "not_ready", "storage is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"service": "ft-music-api", "status": "ready"})
}

func (s *Server) listLibrary(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if utf8.RuneCountInString(query) > maxQuerySize {
		writeError(w, http.StatusBadRequest, "invalid_query", "q must be 256 characters or fewer")
		return
	}
	limit, err := parseBoundedInt(r.URL.Query().Get("limit"), defaultLimit, 1, maxLimit)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_limit", "limit must be an integer between 1 and 100")
		return
	}
	offset, err := parseBoundedInt(r.URL.Query().Get("offset"), 0, 0, int(^uint(0)>>1))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_offset", "offset must be a non-negative integer")
		return
	}
	tracks, err := s.music.Library(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_unavailable", "library is unavailable")
		return
	}
	total := len(tracks)
	start := min(offset, total)
	end := min(start+limit, total)
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-Limit", strconv.Itoa(limit))
	w.Header().Set("X-Offset", strconv.Itoa(offset))
	writeJSON(w, http.StatusOK, tracks[start:end])
}

func (s *Server) getTrack(w http.ResponseWriter, r *http.Request) {
	track, err := s.music.Track(r.Context(), r.PathValue("id"))
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "track_not_found", "track was not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "track_unavailable", "track is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, track)
}

func (s *Server) addTrack(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json")
		return
	}
	var track domain.Track
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxJSONBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&track); err != nil {
		writeDecodeError(w, err)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must contain exactly one JSON object")
		return
	}
	added, err := s.music.AddTrack(r.Context(), track)
	if errors.Is(err, domain.ErrInvalidTrack) {
		writeError(w, http.StatusBadRequest, "invalid_track", strings.TrimPrefix(err.Error(), domain.ErrInvalidTrack.Error()+": "))
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "track_save_failed", "track could not be saved")
		return
	}
	writeJSON(w, http.StatusCreated, added)
}

func (s *Server) deleteTrack(w http.ResponseWriter, r *http.Request) {
	err := s.music.DeleteTrack(r.Context(), r.PathValue("id"))
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "track_not_found", "track was not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "track_delete_failed", "track could not be deleted")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" || utf8.RuneCountInString(query) > maxQuerySize {
		writeError(w, http.StatusBadRequest, "invalid_query", "q is required and must be 256 characters or fewer")
		return
	}
	tracks, err := s.music.SearchRemote(r.Context(), query)
	if errors.Is(err, streaming.ErrProviderUnavailable) {
		writeError(w, http.StatusServiceUnavailable, "providers_unavailable", "no streaming providers are configured")
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		writeError(w, http.StatusGatewayTimeout, "search_timeout", "remote search timed out")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadGateway, "search_failed", "remote search could not be completed")
		return
	}
	writeJSON(w, http.StatusOK, tracks)
}

func (s *Server) listProviders(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"providers": s.music.Providers()})
}

func (s *Server) streamURL(w http.ResponseWriter, r *http.Request) {
	streamURL, err := s.music.StreamURL(r.Context(), r.PathValue("provider"), r.PathValue("id"))
	if errors.Is(err, streaming.ErrProviderNotFound) {
		writeError(w, http.StatusNotFound, "provider_not_found", "streaming provider was not found")
		return
	}
	if errors.Is(err, streaming.ErrInvalidStreamURL) {
		writeError(w, http.StatusBadGateway, "invalid_stream_url", "provider returned an invalid stream URL")
		return
	}
	if errors.Is(err, streaming.ErrProviderUnavailable) {
		writeError(w, http.StatusServiceUnavailable, "provider_unavailable", "streaming provider is unavailable")
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		writeError(w, http.StatusGatewayTimeout, "stream_timeout", "stream URL resolution timed out")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadGateway, "stream_resolution_failed", "stream URL could not be resolved")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": streamURL})
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body exceeds 1 MiB")
		return
	}
	writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON for a track")
}

func parseBoundedInt(raw string, fallback, minValue, maxValue int) (int, error) {
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < minValue || value > maxValue {
		return 0, fmt.Errorf("value outside allowed range")
	}
	return value, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.Header.Get("X-Request-ID"))
		validID := id != "" && len(id) <= 128
		for _, ch := range id {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-' || ch == '_' || ch == '.' || ch == ':') {
				validID = false
				break
			}
		}
		if !validID {
			var random [16]byte
			if _, err := rand.Read(random[:]); err == nil {
				id = hex.EncodeToString(random[:])
			} else {
				id = "ftmusic-request"
			}
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" || r.URL.Path == "/readyz" || r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}
		authorization := strings.TrimSpace(r.Header.Get("Authorization"))
		provided := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		if !strings.HasPrefix(authorization, "Bearer ") || len(provided) != len(s.apiKey) || subtle.ConstantTimeCompare([]byte(provided), []byte(s.apiKey)) != 1 {
			w.Header().Set("WWW-Authenticate", `Bearer realm="FT-music API"`)
			writeError(w, http.StatusUnauthorized, "unauthorized", "a valid bearer token is required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				log.Printf("FT-music API recovered from an unexpected panic")
				writeError(w, http.StatusInternalServerError, "internal_error", "an unexpected server error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
