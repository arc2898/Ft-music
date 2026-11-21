// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package repository

import (
	"context"
	"errors"

	"github.com/arc2898/ft-music/go/internal/domain"
)

var ErrNotFound = errors.New("track not found")

type TrackRepository interface {
	Add(context.Context, domain.Track) (domain.Track, error)
	List(context.Context) ([]domain.Track, error)
	Search(context.Context, string) ([]domain.Track, error)
	Get(context.Context, string) (domain.Track, error)
	Delete(context.Context, string) error
	Ping(context.Context) error
}
