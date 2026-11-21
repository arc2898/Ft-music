// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
)

type Tracks struct {
	mu     sync.RWMutex
	tracks map[string]domain.Track
}

func NewTracks() *Tracks { return &Tracks{tracks: make(map[string]domain.Track)} }

func (r *Tracks) Add(ctx context.Context, t domain.Track) (domain.Track, error) {
	if err := ctx.Err(); err != nil {
		return domain.Track{}, err
	}
	t = t.Normalize()
	if err := t.Validate(); err != nil {
		return domain.Track{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.tracks[t.ID]; ok {
		t.CreatedAt = existing.CreatedAt
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	} else {
		t.CreatedAt = t.CreatedAt.UTC()
	}
	r.tracks[t.ID] = t
	return t, nil
}

func (r *Tracks) List(ctx context.Context) ([]domain.Track, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Track, 0, len(r.tracks))
	for _, t := range r.tracks {
		out = append(out, t)
	}
	sortTracks(out)
	return out, nil
}

func (r *Tracks) Search(ctx context.Context, query string) ([]domain.Track, error) {
	all, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return all, nil
	}
	out := make([]domain.Track, 0)
	for _, t := range all {
		if strings.Contains(strings.ToLower(t.Title), q) || strings.Contains(strings.ToLower(t.Artist), q) || strings.Contains(strings.ToLower(t.Album), q) {
			out = append(out, t)
		}
	}
	return out, nil
}

func (r *Tracks) Get(ctx context.Context, id string) (domain.Track, error) {
	if err := ctx.Err(); err != nil {
		return domain.Track{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tracks[id]
	if !ok {
		return domain.Track{}, repository.ErrNotFound
	}
	return t, nil
}

func (r *Tracks) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tracks[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.tracks, id)
	return nil
}

func (r *Tracks) Ping(ctx context.Context) error { return ctx.Err() }

func sortTracks(tracks []domain.Track) {
	sort.Slice(tracks, func(i, j int) bool {
		if tracks[i].CreatedAt.Equal(tracks[j].CreatedAt) {
			return tracks[i].ID < tracks[j].ID
		}
		return tracks[i].CreatedAt.Before(tracks[j].CreatedAt)
	})
}
