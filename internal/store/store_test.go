package store

import (
	"path/filepath"
	"testing"
)

func TestStoreUpsertAndList(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "routes.json")
	s := NewStore(tmp)

	if err := s.Upsert("meuapp", Route{Port: 3000}); err != nil {
		t.Fatalf("Upsert returned unexpected error: %v", err)
	}

	routes := s.List()

	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}

	got, ok := routes["meuapp"]
	if !ok {
		t.Fatal("expected route 'meuapp' to be present")
	}
	if got.Port != 3000 {
		t.Errorf("expected Port=3000, got Port=%d", got.Port)
	}
}

func TestStoreDelete(t *testing.T) {
	cases := []struct {
		name          string
		seed          map[string]Route
		target        string
		wantErr       bool
		wantRemaining int
	}{
		{
			name:          "removes an existing route",
			seed:          map[string]Route{"foo": {Port: 3000}},
			target:        "foo",
			wantErr:       false,
			wantRemaining: 0,
		},
		{
			name:          "returns error when route is missing",
			seed:          map[string]Route{"foo": {Port: 3000}},
			target:        "bar",
			wantErr:       true,
			wantRemaining: 1,
		},
		{
			name:          "returns error on empty store",
			seed:          nil,
			target:        "anything",
			wantErr:       true,
			wantRemaining: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmp := filepath.Join(t.TempDir(), "routes.json")
			s := NewStore(tmp)
			for name, route := range tc.seed {
				if err := s.Upsert(name, route); err != nil {
					t.Fatalf("seed Upsert failed: %v", err)
				}
			}

			err := s.Delete(tc.target)

			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got := len(s.List()); got != tc.wantRemaining {
				t.Errorf("expected %d remaining routes, got %d", tc.wantRemaining, got)
			}
		})
	}
}
