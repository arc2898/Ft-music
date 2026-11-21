package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
)

func TestTracksAddSearchAndDelete(t *testing.T) {
	repo := NewTracks()
	ctx := context.Background()
	if _, err := repo.Add(ctx, domain.Track{ID: " t1 ", Title: " Signal ", Artist: " FT "}); err != nil {
		t.Fatal(err)
	}
	got, err := repo.Search(ctx, "signal")
	if err != nil || len(got) != 1 || got[0].ID != "t1" {
		t.Fatalf("Search() = %#v, err = %v", got, err)
	}
	if err := repo.Delete(ctx, "missing"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Delete(missing) error = %v", err)
	}
	if err := repo.Delete(ctx, "t1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Get(ctx, "t1"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Get(deleted) error = %v", err)
	}
}
