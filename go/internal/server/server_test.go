package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arc2898/ft-music/go/internal/domain"
	"github.com/arc2898/ft-music/go/internal/repository/memory"
	"github.com/arc2898/ft-music/go/internal/service"
	"github.com/arc2898/ft-music/go/internal/streaming"
)

type testProvider struct{}

func (testProvider) Name() string { return "demo" }
func (testProvider) Search(_ context.Context, query string) ([]domain.Track, error) {
	return []domain.Track{{ID: "remote-1", Title: query, Artist: "FT-music"}}, nil
}
func (testProvider) StreamURL(context.Context, string) (string, error) {
	return "https://media.example.test/track.mp3", nil
}

func testHandler() http.Handler {
	return New(service.NewMusic(memory.NewTracks(), streaming.NewCatalog())).Routes()
}

func TestHealthAndReadiness(t *testing.T) {
	handler := testHandler()
	for _, path := range []string{"/healthz", "/api/v1/health", "/readyz"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", path, response.Code, http.StatusOK)
		}
	}
}

func TestLibraryCRUDAndPagination(t *testing.T) {
	handler := testHandler()
	for _, id := range []string{"t1", "t2", "t3"} {
		body := []byte(`{"id":"` + id + `","title":"Signal ` + id + `","artist":"FT"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/v1/library", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusCreated {
			t.Fatalf("add %s status = %d: %s", id, response.Code, response.Body.String())
		}
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/library?q=signal&limit=2&offset=1", nil))
	if response.Code != http.StatusOK || response.Header().Get("X-Total-Count") != "3" {
		t.Fatalf("list response = %d %s headers=%v", response.Code, response.Body.String(), response.Header())
	}
	var items []domain.Track
	if err := json.Unmarshal(response.Body.Bytes(), &items); err != nil || len(items) != 2 {
		t.Fatalf("list body = %s, err = %v", response.Body.String(), err)
	}

	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/v1/library/t2", nil))
	if get.Code != http.StatusOK || !bytes.Contains(get.Body.Bytes(), []byte(`"id":"t2"`)) {
		t.Fatalf("get response = %d %s", get.Code, get.Body.String())
	}

	remove := httptest.NewRecorder()
	handler.ServeHTTP(remove, httptest.NewRequest(http.MethodDelete, "/api/v1/library/t2", nil))
	if remove.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, body=%s", remove.Code, remove.Body.String())
	}
}

func TestAddTrackRejectsBadRequests(t *testing.T) {
	handler := testHandler()
	cases := []struct {
		name        string
		contentType string
		body        string
		wantStatus  int
	}{
		{"wrong content type", "text/plain", `{}`, http.StatusUnsupportedMediaType},
		{"invalid JSON subtype", "application/jsonp", `{}`, http.StatusUnsupportedMediaType},
		{"invalid JSON", "application/json", `{`, http.StatusBadRequest},
		{"unknown field", "application/json", `{"id":"1","title":"x","artist":"y","unexpected":true}`, http.StatusBadRequest},
		{"trailing JSON", "application/json", `{"id":"1","title":"x","artist":"y"} {}`, http.StatusBadRequest},
		{"invalid track", "application/json", `{"id":"1","title":"x","artist":"y","duration_seconds":-5}`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/library", bytes.NewBufferString(tc.body))
			request.Header.Set("Content-Type", tc.contentType)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.wantStatus, response.Body.String())
			}
		})
	}
}

func TestRemoteSearchAndSafeStreamResolution(t *testing.T) {
	catalog := streaming.NewCatalog(testProvider{})
	handler := New(service.NewMusic(memory.NewTracks(), catalog)).Routes()

	search := httptest.NewRecorder()
	handler.ServeHTTP(search, httptest.NewRequest(http.MethodGet, "/api/v1/search?q=aurora", nil))
	if search.Code != http.StatusOK || !bytes.Contains(search.Body.Bytes(), []byte(`"title":"aurora"`)) {
		t.Fatalf("search response = %d %s", search.Code, search.Body.String())
	}

	stream := httptest.NewRecorder()
	handler.ServeHTTP(stream, httptest.NewRequest(http.MethodGet, "/api/v1/providers/demo/tracks/remote-1/stream", nil))
	if stream.Code != http.StatusOK || !bytes.Contains(stream.Body.Bytes(), []byte("https://media.example.test/track.mp3")) {
		t.Fatalf("stream response = %d %s", stream.Code, stream.Body.String())
	}
}

func TestRequestAndSecurityHeaders(t *testing.T) {
	response := httptest.NewRecorder()
	testHandler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	for key, want := range map[string]string{
		"X-Request-ID":           "",
		"X-Content-Type-Options": "nosniff",
		"Referrer-Policy":        "no-referrer",
		"Cache-Control":          "no-store",
	} {
		got := response.Header().Get(key)
		if got == "" || (want != "" && got != want) {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestAPIKeyProtectsApplicationRoutesButNotHealth(t *testing.T) {
	music := service.NewMusic(memory.NewTracks(), streaming.NewCatalog())
	handler := New(music, WithAPIKey("test-secret")).Routes()

	health := httptest.NewRecorder()
	handler.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", health.Code, http.StatusOK)
	}

	unauthorized := httptest.NewRecorder()
	handler.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/api/v1/library", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}

	authorizedRequest := httptest.NewRequest(http.MethodGet, "/api/v1/library", nil)
	authorizedRequest.Header.Set("Authorization", "Bearer test-secret")
	authorized := httptest.NewRecorder()
	handler.ServeHTTP(authorized, authorizedRequest)
	if authorized.Code != http.StatusOK {
		t.Fatalf("authorized status = %d, want %d", authorized.Code, http.StatusOK)
	}
}

func TestRemoteSearchWithoutProvidersReturnsServiceUnavailable(t *testing.T) {
	response := httptest.NewRecorder()
	testHandler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/search?q=signal", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusServiceUnavailable, response.Body.String())
	}
}
