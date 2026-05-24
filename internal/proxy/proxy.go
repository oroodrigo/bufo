package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/oroodrigo/bufo/internal/store"
)

type Proxy struct {
	store *store.Store
}

func New(s *store.Store) *Proxy {
	return &Proxy{store: s}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := extractName(r.Host)
	if name == "" {
		http.Error(w, "host inválido", http.StatusBadRequest)
		return
	}

	routes := p.store.List()
	route, ok := routes[name]
	if !ok {
		http.Error(w, fmt.Sprintf("rota '%s' não registrada", name), http.StatusNotFound)
		return
	}

	target := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%d", route.Port),
	}
	httputil.NewSingleHostReverseProxy(target).ServeHTTP(w, r)
}

// extractName turns "meuapp.localhost:1355" into "meuapp" and
// "api.frontend.localhost" into "api.frontend". An empty string
// signals the Host header doesn't match the expected pattern.
func extractName(host string) string {
	if i := strings.IndexByte(host, ':'); i >= 0 {
		host = host[:i]
	}

	trimmed := strings.TrimSuffix(host, ".localhost")
	if trimmed == host {
		return ""
	}
	return trimmed
}
