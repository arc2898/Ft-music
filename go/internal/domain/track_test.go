package domain

import (
	"errors"
	"testing"
)

func TestTrackNormalizeAndValidate(t *testing.T) {
	track := (Track{ID: " id ", Title: " Signal ", Artist: " FT "}).Normalize()
	if track.ID != "id" || track.Title != "Signal" || track.Artist != "FT" {
		t.Fatalf("Normalize() = %#v", track)
	}
	if err := track.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestTrackRejectsInvalidMetadata(t *testing.T) {
	cases := []Track{
		{ID: "1", Artist: "FT"},
		{ID: "1", Title: "Signal"},
		{ID: "1", Title: "Signal", Artist: "FT", Duration: -1},
		{ID: "1", Title: "Signal", Artist: "FT", ArtworkURL: "file:///tmp/a.png"},
	}
	for i, track := range cases {
		if err := track.Validate(); !errors.Is(err, ErrInvalidTrack) {
			t.Errorf("case %d: Validate() error = %v, want ErrInvalidTrack", i, err)
		}
	}
}
