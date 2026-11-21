// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS tracks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    album TEXT NOT NULL DEFAULT '',
    duration_seconds INTEGER NOT NULL DEFAULT 0 CHECK (duration_seconds >= 0),
    artwork_url TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_tracks_title ON tracks(title);
CREATE INDEX IF NOT EXISTS idx_tracks_artist ON tracks(artist);
CREATE INDEX IF NOT EXISTS idx_tracks_album ON tracks(album);
`

const trackColumns = `id, title, artist, album, duration_seconds, artwork_url, created_at`

type Tracks struct{ db *sql.DB }

func Open(path string) (*Tracks, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("database path cannot be empty")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open SQLite database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping SQLite database: %w", err)
	}
	for _, pragma := range []string{
		`PRAGMA foreign_keys = ON`,
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA journal_mode = WAL`,
		`PRAGMA synchronous = NORMAL`,
	} {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("configure SQLite (%s): %w", pragma, err)
		}
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize SQLite schema: %w", err)
	}
	return &Tracks{db: db}, nil
}

func (r *Tracks) Close() error { return r.db.Close() }

func (r *Tracks) Ping(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return fmt.Errorf("SQLite health check: %w", err)
	}
	return nil
}

func (r *Tracks) Add(ctx context.Context, t domain.Track) (domain.Track, error) {
	t = t.Normalize()
	if err := t.Validate(); err != nil {
		return domain.Track{}, err
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	} else {
		t.CreatedAt = t.CreatedAt.UTC()
	}
	_, err := r.db.ExecContext(ctx, `
INSERT INTO tracks (id, title, artist, album, duration_seconds, artwork_url, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    title = excluded.title,
    artist = excluded.artist,
    album = excluded.album,
    duration_seconds = excluded.duration_seconds,
    artwork_url = excluded.artwork_url`,
		t.ID, t.Title, t.Artist, t.Album, t.Duration, t.ArtworkURL, t.CreatedAt.Format(time.RFC3339Nano))
	if err != nil {
		return domain.Track{}, fmt.Errorf("save track %q: %w", t.ID, err)
	}
	return r.Get(ctx, t.ID)
}

func (r *Tracks) List(ctx context.Context) ([]domain.Track, error) {
	return r.query(ctx, `SELECT `+trackColumns+` FROM tracks ORDER BY created_at ASC, id ASC`)
}

func (r *Tracks) Search(ctx context.Context, query string) ([]domain.Track, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return r.List(ctx)
	}
	q = escapeLike(q)
	pattern := "%" + q + "%"
	return r.query(ctx, `SELECT `+trackColumns+` FROM tracks
WHERE title LIKE ? ESCAPE '\' OR artist LIKE ? ESCAPE '\' OR album LIKE ? ESCAPE '\'
ORDER BY created_at ASC, id ASC`, pattern, pattern, pattern)
}

func (r *Tracks) Get(ctx context.Context, id string) (domain.Track, error) {
	t, err := scanTrack(r.db.QueryRowContext(ctx, `SELECT `+trackColumns+` FROM tracks WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Track{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Track{}, fmt.Errorf("get track %q: %w", id, err)
	}
	return t, nil
}

func (r *Tracks) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tracks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete track %q: %w", id, err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted track %q: %w", id, err)
	}
	if count == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *Tracks) query(ctx context.Context, query string, args ...any) ([]domain.Track, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query tracks: %w", err)
	}
	defer rows.Close()

	tracks := make([]domain.Track, 0)
	for rows.Next() {
		t, err := scanTrack(rows)
		if err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}
		tracks = append(tracks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tracks: %w", err)
	}
	return tracks, nil
}

type scanner interface{ Scan(...any) error }

func scanTrack(row scanner) (domain.Track, error) {
	var t domain.Track
	var created string
	if err := row.Scan(&t.ID, &t.Title, &t.Artist, &t.Album, &t.Duration, &t.ArtworkURL, &created); err != nil {
		return domain.Track{}, err
	}
	parsed, err := time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return domain.Track{}, fmt.Errorf("parse created_at: %w", err)
	}
	t.CreatedAt = parsed.UTC()
	return t, nil
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(``, `\`, "%", `\%`, "_", `\_`)
	return replacer.Replace(value)
}
