// FT-music is an open-source music application.
// Copyright (C) 2026 FT-music contributors
// SPDX-License-Identifier: GPL-2.0-or-later

package streaming

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/arc2898/ft-music/go/internal/domain"
)

var (
	ErrProviderUnavailable = errors.New("streaming provider unavailable")
	ErrProviderNotFound    = errors.New("streaming provider not found")
	ErrInvalidProvider     = errors.New("invalid streaming provider")
	ErrInvalidStreamURL    = errors.New("provider returned an invalid stream URL")
)

const (
	defaultSearchTimeout = 12 * time.Second
	searchCacheTTL        = 2 * time.Minute
	searchCacheCapacity   = 128
)

var providerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,79}$`)

// Provider is the only boundary the application uses for remote music.
// Implementations must use documented, authorized provider APIs and respect
// their terms, licenses, rate limits, attribution, and regional restrictions.
type Provider interface {
	Name() string
	Search(context.Context, string) ([]domain.Track, error)
	StreamURL(context.Context, string) (string, error)
}

type searchCacheEntry struct {
	key       string
	tracks    []domain.Track
	expiresAt time.Time
}

type searchCall struct {
	done   chan struct{}
	tracks []domain.Track
	err    error
}

type Catalog struct {
	mu            sync.Mutex
	providers     map[string]Provider
	searchTimeout time.Duration
	cache         map[string]*list.Element
	cacheOrder    *list.List
	inFlight      map[string]*searchCall
}

func NewCatalog(providers ...Provider) *Catalog {
	catalog := &Catalog{
		providers:     make(map[string]Provider, len(providers)),
		searchTimeout: defaultSearchTimeout,
		cache:         make(map[string]*list.Element),
		cacheOrder:    list.New(),
		inFlight:      make(map[string]*searchCall),
	}
	for _, provider := range providers {
		_ = catalog.Register(provider)
	}
	return catalog
}

func (c *Catalog) Register(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("%w: provider is nil", ErrInvalidProvider)
	}
	name := strings.TrimSpace(provider.Name())
	if !providerNamePattern.MatchString(name) {
		return fmt.Errorf("%w: provider name must be a 1 to 80 character slug", ErrInvalidProvider)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.providers[name]; exists {
		return fmt.Errorf("%w: %s is already registered", ErrInvalidProvider, name)
	}
	c.providers[name] = provider
	// A newly registered provider changes the result set for every query.
	c.cache = make(map[string]*list.Element)
	c.cacheOrder.Init()
	return nil
}

func (c *Catalog) Providers() []Provider {
	c.mu.Lock()
	providers := make([]Provider, 0, len(c.providers))
	for _, provider := range c.providers {
		providers = append(providers, provider)
	}
	c.mu.Unlock()
	sort.Slice(providers, func(i, j int) bool { return providers[i].Name() < providers[j].Name() })
	return providers
}

func (c *Catalog) ProviderNames() []string {
	providers := c.Providers()
	names := make([]string, 0, len(providers))
	for _, provider := range providers {
		names = append(names, provider.Name())
	}
	return names
}

// Search uses a bounded TTL/LRU cache and coalesces concurrent requests for the
// same normalized query. The cached values are copied on read/write so callers
// cannot mutate a result shared with other requests.
func (c *Catalog) Search(ctx context.Context, query string) ([]domain.Track, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	key := strings.ToLower(query)

	c.mu.Lock()
	if cached := c.getCachedLocked(key, time.Now()); cached != nil {
		c.mu.Unlock()
		return cloneTracks(cached), nil
	}
	if call, exists := c.inFlight[key]; exists {
		c.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-call.done:
			return cloneTracks(call.tracks), call.err
		}
	}
	call := &searchCall{done: make(chan struct{})}
	c.inFlight[key] = call
	c.mu.Unlock()

	tracks, err := c.searchProviders(ctx, query)

	c.mu.Lock()
	call.tracks = cloneTracks(tracks)
	call.err = err
	if err == nil {
		c.putCachedLocked(key, tracks, time.Now().Add(searchCacheTTL))
	}
	delete(c.inFlight, key)
	close(call.done)
	c.mu.Unlock()

	return cloneTracks(tracks), err
}

func (c *Catalog) searchProviders(ctx context.Context, query string) ([]domain.Track, error) {
	providers := c.Providers()
	if len(providers) == 0 {
		return nil, ErrProviderUnavailable
	}

	searchCtx, cancel := context.WithTimeout(ctx, c.searchTimeout)
	defer cancel()
	type result struct {
		index  int
		tracks []domain.Track
		err    error
	}
	results := make(chan result, len(providers))
	for i, provider := range providers {
		go func(index int, p Provider) {
			tracks, err := safeProviderSearch(searchCtx, p, query)
			results <- result{index: index, tracks: tracks, err: err}
		}(i, provider)
	}

	ordered := make([]result, len(providers))
	for received := 0; received < len(providers); received++ {
		select {
		case item := <-results:
			ordered[item.index] = item
		case <-searchCtx.Done():
			received = len(providers)
		}
	}

	tracks := make([]domain.Track, 0)
	failures := make([]error, 0)
	seen := make(map[string]int)
	for _, item := range ordered {
		if item.err != nil {
			failures = append(failures, item.err)
			continue
		}
		for _, track := range item.tracks {
			track = track.Normalize()
			if track.Validate() != nil {
				continue
			}
			key := strings.ToLower(track.Title) + "\x00" + strings.ToLower(track.Artist)
			if existingIndex, ok := seen[key]; ok {
				if tracks[existingIndex].Album == "" && track.Album != "" {
					tracks[existingIndex].Album = track.Album
				}
				if tracks[existingIndex].ArtworkURL == "" && track.ArtworkURL != "" {
					tracks[existingIndex].ArtworkURL = track.ArtworkURL
				}
				continue
			}
			seen[key] = len(tracks)
			tracks = append(tracks, track)
		}
	}
	if len(tracks) > 0 {
		return tracks, nil
	}
	if searchCtx.Err() != nil {
		return nil, searchCtx.Err()
	}
	if len(failures) > 0 {
		return nil, errors.Join(failures...)
	}
	return tracks, nil
}

func (c *Catalog) getCachedLocked(key string, now time.Time) []domain.Track {
	element, ok := c.cache[key]
	if !ok {
		return nil
	}
	entry := element.Value.(*searchCacheEntry)
	if !now.Before(entry.expiresAt) {
		delete(c.cache, key)
		c.cacheOrder.Remove(element)
		return nil
	}
	c.cacheOrder.MoveToFront(element)
	return cloneTracks(entry.tracks)
}

func (c *Catalog) putCachedLocked(key string, tracks []domain.Track, expiresAt time.Time) {
	if existing, ok := c.cache[key]; ok {
		c.cacheOrder.Remove(existing)
	}
	element := c.cacheOrder.PushFront(&searchCacheEntry{
		key: key, tracks: cloneTracks(tracks), expiresAt: expiresAt,
	})
	c.cache[key] = element
	for c.cacheOrder.Len() > searchCacheCapacity {
		oldest := c.cacheOrder.Back()
		entry := oldest.Value.(*searchCacheEntry)
		delete(c.cache, entry.key)
		c.cacheOrder.Remove(oldest)
	}
}

func cloneTracks(tracks []domain.Track) []domain.Track {
	if tracks == nil {
		return nil
	}
	cloned := make([]domain.Track, len(tracks))
	copy(cloned, tracks)
	return cloned
}

func (c *Catalog) StreamURL(ctx context.Context, providerName, trackID string) (string, error) {
	if strings.TrimSpace(trackID) == "" {
		return "", errors.New("track id cannot be empty")
	}
	c.mu.Lock()
	provider, ok := c.providers[providerName]
	c.mu.Unlock()
	if !ok {
		return "", ErrProviderNotFound
	}
	streamCtx, cancel := context.WithTimeout(ctx, defaultSearchTimeout)
	defer cancel()
	streamURL, err := safeProviderStreamURL(streamCtx, provider, trackID)
	if err != nil {
		return "", err
	}
	parsed, err := url.ParseRequestURI(strings.TrimSpace(streamURL))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil {
		return "", ErrInvalidStreamURL
	}
	return parsed.String(), nil
}

func safeProviderSearch(ctx context.Context, provider Provider, query string) (tracks []domain.Track, err error) {
	defer func() {
		if recover() != nil {
			tracks = nil
			err = errors.New("provider search panicked")
		}
	}()
	return provider.Search(ctx, query)
}

func safeProviderStreamURL(ctx context.Context, provider Provider, id string) (streamURL string, err error) {
	defer func() {
		if recover() != nil {
			streamURL = ""
			err = errors.New("provider stream resolution panicked")
		}
	}()
	return provider.StreamURL(ctx, id)
}
