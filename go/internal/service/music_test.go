package service

import (
	"context"
	"errors"
	"testing"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository/memory"
	"github.com/arc2898/ft-music/go/internal/streaming"
)

func TestLibraryDelegatesToRepository(t *testing.T) {
	music := NewMusic(memory.NewTracks(), streaming.NewCatalog())
	ctx := context.Background()
	if _, err := music.AddTrack(ctx, domain.Track{ID: "1", Title: "Midnight Drive", Artist: "FT"}); err != nil {
		t.Fatal(err)
	}
	got, err := music.Library(ctx, "midnight")
	if err != nil || len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("Library() = %#v, err=%v", got, err)
	}
}

func TestAddTrackValidatesDomain(t *testing.T) {
	music := NewMusic(memory.NewTracks(), streaming.NewCatalog())
	_, err := music.AddTrack(context.Background(), domain.Track{ID: "", Title: "No ID", Artist: "FT"})
	if !errors.Is(err, domain.ErrInvalidTrack) {
		t.Fatalf("AddTrack() error = %v", err)
	}
}
