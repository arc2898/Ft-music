package sqlite

import (
	"context"
	"errors"
	"testing"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository"
)

func TestTracksPersistSearchAndDelete(t *testing.T) {
	r, err := Open("file:ft_music_test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	ctx := context.Background()
	if _, err := r.Add(ctx, domain.Track{ID: "1", Title: "Signal_One", Artist: "FT"}); err != nil {
		t.Fatal(err)
	}
	got, err := r.Search(ctx, "Signal_")
	if err != nil || len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("Search() = %#v, err = %v", got, err)
	}
	if _, err := r.Get(ctx, "missing"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Get(missing) error = %v", err)
	}
	if err := r.Delete(ctx, "1"); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(ctx, "1"); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Delete(missing) error = %v", err)
	}
}
