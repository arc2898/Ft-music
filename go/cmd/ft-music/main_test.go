package main

import "testing"

func TestNonLoopbackAddressRequiresAuthentication(t *testing.T) {
	cases := []struct {
		address string
		want    bool
	}{
		{"127.0.0.1:8080", false},
		{"localhost:8080", false},
		{"[::1]:8080", false},
		{":8080", true},
		{"0.0.0.0:8080", true},
		{"192.0.2.10:8080", true},
		{"music.example.test:8080", true},
	}
	for _, tc := range cases {
		if got := isNonLoopbackAddress(tc.address); got != tc.want {
			t.Errorf("isNonLoopbackAddress(%q) = %v, want %v", tc.address, got, tc.want)
		}
	}
}
