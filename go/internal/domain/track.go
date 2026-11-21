// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package domain

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidTrack = errors.New("invalid track")

type Track struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Album      string    `json:"album,omitempty"`
	Duration   int       `json:"duration_seconds,omitempty"`
	ArtworkURL string    `json:"artwork_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Normalize trims user- and provider-supplied metadata without changing its meaning.
func (t Track) Normalize() Track {
	t.ID = strings.TrimSpace(t.ID)
	t.Title = strings.TrimSpace(t.Title)
	t.Artist = strings.TrimSpace(t.Artist)
	t.Album = strings.TrimSpace(t.Album)
	t.ArtworkURL = strings.TrimSpace(t.ArtworkURL)
	return t
}

// Validate enforces the stable contract shared by storage and API clients.
func (t Track) Validate() error {
	checks := []struct {
		name  string
		value string
		max   int
	}{
		{"id", t.ID, 256},
		{"title", t.Title, 500},
		{"artist", t.Artist, 500},
		{"album", t.Album, 500},
		{"artwork_url", t.ArtworkURL, 2048},
	}
	for _, check := range checks {
		if check.name == "id" || check.name == "title" || check.name == "artist" {
			if strings.TrimSpace(check.value) == "" {
				return fmt.Errorf("%w: %s is required", ErrInvalidTrack, check.name)
			}
		}
		if utf8.RuneCountInString(check.value) > check.max {
			return fmt.Errorf("%w: %s exceeds %d characters", ErrInvalidTrack, check.name, check.max)
		}
		if strings.IndexFunc(check.value, unicode.IsControl) >= 0 {
			return fmt.Errorf("%w: %s contains control characters", ErrInvalidTrack, check.name)
		}
	}
	if t.Duration < 0 {
		return fmt.Errorf("%w: duration_seconds cannot be negative", ErrInvalidTrack)
	}
	if t.ArtworkURL != "" {
		parsed, err := url.ParseRequestURI(t.ArtworkURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil {
			return fmt.Errorf("%w: artwork_url must be an absolute HTTP or HTTPS URL", ErrInvalidTrack)
		}
	}
	return nil
}
