package proxy

import "testing"

func TestExtractName(t *testing.T) {
	cases := []struct {
		name string
		host string
		want string
	}{
		{"simple host with port", "meuapp.localhost:1355", "meuapp"},
		{"simple host without port", "meuapp.localhost", "meuapp"},
		{"nested subdomain", "api.frontend.localhost:1355", "api.frontend"},
		{"unrelated host", "example.com:80", ""},
		{"plain localhost", "localhost:1355", ""},
		{"empty host", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractName(tc.host)
			if got != tc.want {
				t.Errorf("extractName(%q) = %q, want %q", tc.host, got, tc.want)
			}
		})
	}
}
