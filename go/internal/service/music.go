// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package service

import (
	"context"
	"errors"
	"strings"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
	"github.com/arc2898/ft-music/go/internal/streaming"
)

var ErrInvalidQuery = errors.New("search query cannot be empty")

type Music struct {
	tracks    repository.TrackRepository
	streaming *streaming.Catalog
}

func NewMusic(tracks repository.TrackRepository, catalog *streaming.Catalog) *Music {
	return &Music{tracks: tracks, streaming: catalog}
}

func (m *Music) AddTrack(ctx context.Context, track domain.Track) (domain.Track, error) {
	track = track.Normalize()
	if err := track.Validate(); err != nil {
		return domain.Track{}, err
	}
	return m.tracks.Add(ctx, track)
}

func (m *Music) Library(ctx context.Context, query string) ([]domain.Track, error) {
	return m.tracks.Search(ctx, strings.TrimSpace(query))
}

func (m *Music) Track(ctx context.Context, id string) (domain.Track, error) {
	return m.tracks.Get(ctx, strings.TrimSpace(id))
}

func (m *Music) DeleteTrack(ctx context.Context, id string) error {
	return m.tracks.Delete(ctx, strings.TrimSpace(id))
}

func (m *Music) Ready(ctx context.Context) error { return m.tracks.Ping(ctx) }

func (m *Music) Providers() []string { return m.streaming.ProviderNames() }

func (m *Music) SearchRemote(ctx context.Context, query string) ([]domain.Track, error) {
	if strings.TrimSpace(query) == "" {
		return nil, ErrInvalidQuery
	}
	return m.streaming.Search(ctx, query)
}

func (m *Music) StreamURL(ctx context.Context, provider, id string) (string, error) {
	return m.streaming.StreamURL(ctx, strings.TrimSpace(provider), strings.TrimSpace(id))
}
