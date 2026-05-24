package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oroodrigo/bufo/internal/config"
	"github.com/oroodrigo/bufo/internal/store"
)

var (
	s *store.Store
)

func Start() {
	c := config.Load()

	if err := os.MkdirAll(c.BufoDir, 0755); err != nil {
		panic(err)
	}

	if running, pid, _ := Status(c.PIDFile, c.SocketFile); running {
		fmt.Fprintf(os.Stderr, "daemon já está rodando (PID: %d)\n", pid)
		os.Exit(1)
	}

	// Remove leftovers from a previous dirty exit so net.Listen can bind.
	os.Remove(c.SocketFile)
	os.Remove(c.PIDFile)

	if err := savePID(c.PIDFile); err != nil {
		panic(err)
	}

	s = store.NewStore(c.StoreFile)

	cleanup := func() {
		os.Remove(c.SocketFile)
		os.Remove(c.PIDFile)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cleanup()
		os.Exit(0)
	}()

	defer cleanup()
	startServer(c.SocketFile)
}

func Stop(pidFile string, socketFile string) error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("daemon não está rodando")
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("PID inválido")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("processo não encontrado")
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("erro ao matar processo: %w", err)
	}

	waitProcessExit(pid, 2*time.Second)

	os.Remove(pidFile)
	os.Remove(socketFile)
	return nil
}

func waitProcessExit(pid int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return
		}
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func Status(pidFile string, socketFile string) (bool, int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false, 0, nil
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false, 0, fmt.Errorf("PID inválido")
	}

	if _, err := os.Stat(socketFile); os.IsNotExist(err) {
		return false, 0, nil
	}

	conn, err := net.DialTimeout("unix", socketFile, 500*time.Millisecond)
	if err != nil {
		return false, 0, nil
	}
	conn.Close()

	return true, pid, nil
}

func savePID(path string) error {
	pid := os.Getpid()
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func startServer(path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("erro ao remover socket antigo: %w", err))
	}

	l, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Get("/routes", handleList)
	r.Post("/routes/{name}", handleUpsert)
	r.Delete("/routes/{name}", handleDelete)

	if err := http.Serve(l, r); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("servidor encerrou com erro: %w", err))
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	routes := s.List()
	jsonRoutes, err := json.Marshal(routes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRoutes)
}

func handleUpsert(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var route store.Route
	if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
		http.Error(w, "body inválido", http.StatusBadRequest)
		return
	}

	if route.Port <= 0 || route.Port > 65535 {
		http.Error(w, "port inválido (1-65535)", http.StatusBadRequest)
		return
	}

	if err := s.Upsert(name, route); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	if err := s.Delete(name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
