# Bufo 🐸

> Replace `localhost:3000` with clean, readable URLs like `api.myapp.localhost:1355`

Bufo is a local reverse proxy CLI for developers. Instead of juggling port numbers, you give each project a name and access it through a stable `.localhost` URL. Built in Go as a learning project focused on HTTP servers, concurrency, and Unix socket communication.

---

## How it works

Bufo runs a background daemon that listens on port `1355`. When you register a project with `bufo add`, any request to `<name>.localhost:1355` is proxied to the local port your app is running on.

```
Browser → meuapp.localhost:1355 → Bufo daemon → localhost:3000
```

The daemon communicates with the CLI over a Unix socket (`~/.bufo/bufo.sock`), and persists route configuration to `~/.bufo/routes.json`.

---

## Installation

**Requirements:** Go 1.21+

```bash
git clone https://github.com/oroodrigo/bufo
cd bufo/cmd/bufo
go install .
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`.

---

## Usage

```bash
# Add a project
bufo add meuapp --port 3000
# → now accessible at http://meuapp.localhost:1355

# Subdomains work too
bufo add api.myapp --port 10010
# → http://api.myapp.localhost:1355

# List all registered routes
bufo list

# Remove a project
bufo remove meuapp
```

### Daemon management

The daemon starts automatically when you run any `bufo` command. You can also manage it manually:

```bash
bufo daemon start    # start the daemon
bufo daemon stop     # stop the daemon
bufo daemon restart  # restart the daemon
bufo daemon status   # check if the daemon is running
```

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go |
| CLI framework | [Cobra](https://github.com/spf13/cobra) |
| HTTP router | [Chi](https://github.com/go-chi/chi) |
| IPC | Unix socket (`~/.bufo/bufo.sock`) |
| Persistence | JSON (`~/.bufo/routes.json`) |

---

## Architecture

```
cmd/
  bufo/
    main.go           # entrypoint — delegates to internal/cli
internal/
  cli/                # CLI commands (add, remove, list, daemon)
  daemon/             # daemon lifecycle + HTTP API over Unix socket
  store/              # route persistence (routes.json)
  config/             # centralized paths and constants
  proxy/              # reverse proxy (port 1355) — coming soon
```

The CLI is intentionally thin — it validates input, ensures the daemon is running, and forwards requests to it. All business logic lives in the daemon.

---

## Roadmap

- [x] CLI skeleton (add, remove, list, daemon)
- [x] Unix socket communication between CLI and daemon
- [x] Route persistence to `~/.bufo/routes.json`
- [x] Reverse proxy on port 1355
- [ ] Route health polling (every 10s)
- [ ] HTTPS with auto-generated local certificates
- [ ] Clean URLs without port number (`meuapp.localhost`)
- [ ] Process lifecycle management (`bufo run -- npm run dev`)

---

## Platform

Currently developed and tested on **Windows 11**.

---

## Motivation

This project was built as a educational Go project, deliberately chosen to learn HTTP servers, Unix sockets, and concurrency in a practical context. Inspired by [Portless](https://github.com/vercel-labs/portless) from Vercel Labs.
