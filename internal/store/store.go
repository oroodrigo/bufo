package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Route struct {
	Port  int
	Https bool
}

type Store struct {
	path   string
	mu     sync.RWMutex
	routes map[string]Route
}

func NewStore(path string) *Store {
	return &Store{
		path:   path,
		routes: loadStore(path),
	}
}

func loadStore(path string) map[string]Route {
	routes := make(map[string]Route)

	file, err := os.ReadFile(path)
	if err != nil {
		return routes
	}

	if err := json.Unmarshal(file, &routes); err != nil {
		panic(fmt.Errorf("arquivo de rotas %s corrompido: %w", path, err))
	}
	return routes
}

func (s *Store) List() map[string]Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]Route, len(s.routes))
	for k, v := range s.routes {
		copy[k] = v
	}
	return copy
}

func (s *Store) Upsert(routeName string, r Route) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.routes[routeName] = r
	return s.persist()
}

func (s *Store) Delete(routeName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.routes[routeName]
	if !ok {
		return fmt.Errorf("Rota %s não encontrada", routeName)
	}

	delete(s.routes, routeName)
	return s.persist()
}

func (s *Store) persist() error {
	data, err := json.MarshalIndent(s.routes, "", "  ")
	if err != nil {
		return fmt.Errorf("Erro ao serializar rotas: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("Erro ao salvar arquivo: %w", err)
	}

	if err := os.Rename(tmp, s.path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("Erro ao salvar arquivo: %w", err)
	}

	return nil
}
