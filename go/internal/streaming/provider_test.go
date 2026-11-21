package streaming

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/arc2898/ft-music/go/internal/domain"
)

type fakeProvider struct {
	name   string
	tracks []domain.Track
	url    string
	err    error
}

func (p fakeProvider) Name() string { return p.name }
func (p fakeProvider) Search(context.Context, string) ([]domain.Track, error) {
	return p.tracks, p.err
}
func (p fakeProvider) StreamURL(context.Context, string) (string, error) {
	return p.url, p.err
}

func TestCatalogSearchAggregatesAndDeduplicates(t *testing.T) {
	catalog := NewCatalog(
		fakeProvider{name: "second", tracks: []domain.Track{{ID: "2", Title: "Signal", Artist: "FT", Album: "Album"}}},
		fakeProvider{name: "first", tracks: []domain.Track{{ID: "1", Title: " signal ", Artist: "ft"}}},
	)
	got, err := catalog.Search(context.Background(), "signal")
	if err != nil || len(got) != 1 || got[0].Album != "Album" {
		t.Fatalf("Search() = %#v, err = %v", got, err)
	}
}

func TestCatalogRejectsUnsafeStreamURL(t *testing.T) {
	catalog := NewCatalog(fakeProvider{name: "demo", url: "javascript:alert(1)"})
	if _, err := catalog.StreamURL(context.Background(), "demo", "track"); !errors.Is(err, ErrInvalidStreamURL) {
		t.Fatalf("StreamURL() error = %v", err)
	}
}

func TestCatalogReportsMissingProvider(t *testing.T) {
	catalog := NewCatalog()
	if _, err := catalog.Search(context.Background(), "signal"); !errors.Is(err, ErrProviderUnavailable) {
		t.Fatalf("Search() error = %v", err)
	}
}

type panickingProvider struct{}

func (panickingProvider) Name() string { return "panic-provider" }
func (panickingProvider) Search(context.Context, string) ([]domain.Track, error) {
	panic("provider panic must not crash the service")
}
func (panickingProvider) StreamURL(context.Context, string) (string, error) {
	panic("provider panic must not crash the service")
}

func TestCatalogContainsProviderPanics(t *testing.T) {
	catalog := NewCatalog(panickingProvider{})
	if _, err := catalog.Search(context.Background(), "signal"); err == nil {
		t.Fatal("Search() expected a provider error")
	}
	if _, err := catalog.StreamURL(context.Background(), "panic-provider", "track"); err == nil {
		t.Fatal("StreamURL() expected a provider error")
	}
}

type countingProvider struct {
	calls atomic.Int32
}

func (p *countingProvider) Name() string { return "counting" }
func (p *countingProvider) Search(_ context.Context, query string) ([]domain.Track, error) {
	p.calls.Add(1)
	time.Sleep(40 * time.Millisecond)
	return []domain.Track{{ID: "cached-1", Title: query, Artist: "FT-music"}}, nil
}
func (p *countingProvider) StreamURL(context.Context, string) (string, error) {
	return "https://media.example.test/cached.mp3", nil
}

func TestCatalogCoalescesConcurrentSearchesAndCachesResults(t *testing.T) {
	provider := &countingProvider{}
	catalog := NewCatalog(provider)
	const callers = 8

	var wait sync.WaitGroup
	wait.Add(callers)
	errs := make(chan error, callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wait.Done()
			tracks, err := catalog.Search(context.Background(), "same query")
			if err != nil {
				errs <- err
				return
			}
			if len(tracks) != 1 || tracks[0].ID != "cached-1" {
				errs <- errors.New("unexpected search result")
			}
		}()
	}
	wait.Wait()
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	if got := provider.calls.Load(); got != 1 {
		t.Fatalf("provider calls = %d, want 1 for coalesced requests", got)
	}

	first, err := catalog.Search(context.Background(), "same query")
	if err != nil {
		t.Fatal(err)
	}
	first[0].Title = "tampered"
	second, err := catalog.Search(context.Background(), "same query")
	if err != nil {
		t.Fatal(err)
	}
	if second[0].Title != "same query" {
		t.Fatalf("cached result was mutated: %#v", second[0])
	}
	if got := provider.calls.Load(); got != 1 {
		t.Fatalf("provider calls after cache hits = %d, want 1", got)
	}
}

func TestRegisterInvalidatesSearchCache(t *testing.T) {
	first := &countingProvider{}
	catalog := NewCatalog(first)
	if _, err := catalog.Search(context.Background(), "signal"); err != nil {
		t.Fatal(err)
	}

	second := fakeProvider{name: "another", tracks: []domain.Track{{ID: "2", Title: "Signal", Artist: "Other"}}}
	if err := catalog.Register(second); err != nil {
		t.Fatal(err)
	}
	got, err := catalog.Search(context.Background(), "signal")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("Search() after provider registration = %#v, want two tracks", got)
	}
}
