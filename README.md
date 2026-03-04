<h1 align="center">Envio</h1>

<p align="center">
  Local Docker development environments, generated in seconds.
</p>

<p align="center">
  <a href="https://github.com/danjdewhurst/envio/blob/main/LICENSE"><img src="https://img.shields.io/github/license/danjdewhurst/envio" alt="License"></a>
  <a href="https://github.com/danjdewhurst/envio/releases"><img src="https://img.shields.io/github/v/release/danjdewhurst/envio" alt="Latest Release"></a>
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white" alt="Go 1.24">
</p>

---

## Features

- **One command setup** — `envio init laravel` generates a full Docker Compose environment
- **Addon system** — add MySQL, PostgreSQL, Redis, Meilisearch, and more with `--addon`
- **Local `.test` domains** — built-in Traefik reverse proxy routes `*.test` to your containers
- **Trusted HTTPS** — automatic mkcert certificates, no browser warnings
- **App variants** — swap between Nginx + PHP-FPM and FrankenPHP with `--variant`
- **Extensible** — add new apps and addons by implementing a Go interface

---

## Quick Start

```bash
# Install (auto-detects OS and architecture)
curl -sL https://raw.githubusercontent.com/danjdewhurst/envio/main/install.sh | sh

# Create a Laravel project with MySQL and Redis
envio init laravel --addon mysql --addon redis

# Start the environment
envio up
```

Or install via Go:

```bash
go install github.com/danjdewhurst/envio@latest
```

---

## Commands

| Command | Description |
|---------|-------------|
| `envio init <app>` | Generate a Docker Compose environment |
| `envio up` | Start containers (`docker compose up -d`) |
| `envio down` | Stop containers (`docker compose down`) |
| `envio status` | Show config and container status |
| `envio apps` | List available app types |
| `envio addons` | List available addons |
| `envio proxy start` | Start the shared Traefik reverse proxy |
| `envio proxy stop` | Stop the proxy |
| `envio proxy status` | Check proxy status |
| `envio proxy setup-tls` | Install the mkcert root CA for trusted HTTPS |

### `envio init` Flags

```
--variant, -v    App variant (e.g. frankenphp)
--addon, -a      Add a service (repeatable)
--domain, -d     Domain for .test routing (defaults to directory name)
--no-proxy       Disable Traefik proxy integration
```

---

## Supported Apps

| App | Description | Variants |
|-----|-------------|----------|
| `laravel` | PHP Laravel with Nginx + PHP-FPM | `frankenphp` |

---

## Supported Addons

| Addon | Description |
|-------|-------------|
| `mariadb` | MariaDB 11.4 LTS |
| `mysql` | MySQL 8.0 |
| `postgres` | PostgreSQL 16 |
| `redis` | Redis |
| `meilisearch` | Meilisearch |

---

## Local `.test` Domains & HTTPS

Envio uses a shared [Traefik](https://traefik.io/) reverse proxy to route `*.test` domains to your containers, with automatic HTTP → HTTPS redirection.

1. Install [mkcert](https://github.com/FiloSottile/mkcert): `brew install mkcert`
2. Install the local CA: `envio proxy setup-tls`
3. Start the proxy: `envio proxy start`
4. Initialise a project: `envio init laravel --domain myapp && envio up`
5. Visit `https://myapp.test` — no browser warnings.

If mkcert is not installed, everything still works over plain HTTP.

---

## Extending Envio

### Adding a new app

1. Create a package under `internal/app/<name>/`
2. Implement the `app.App` interface:

```go
type App interface {
    Name() string
    DisplayName() string
    Description() string
    Services() []compose.Service
    Volumes() []compose.Volume
    DefaultEnv() map[string]string
    AvailableAddons() []string
}
```

3. Register it in `internal/registry/defaults.go`

### Adding a new addon

1. Create a package under `internal/addon/<name>/`
2. Implement the `addon.Addon` interface:

```go
type Addon interface {
    Name() string
    DisplayName() string
    Description() string
    Services() []compose.Service
    Volumes() []compose.Volume
    EnvVars() map[string]string
}
```

3. Register it in `internal/registry/defaults.go`

---

## Project Structure

```
envio/
├── main.go                            # Entry point
├── cmd/                               # CLI commands (Cobra)
│   ├── root.go                        # Root command + registry init
│   ├── init.go                        # envio init <app> --addon <name>
│   ├── up.go                          # envio up
│   ├── down.go                        # envio down
│   ├── status.go                      # envio status
│   ├── proxy.go                       # envio proxy start/stop/status/setup-*
│   ├── apps.go                        # envio apps
│   └── addons.go                      # envio addons
├── internal/
│   ├── app/                           # App interface + implementations
│   │   ├── app.go
│   │   └── laravel/
│   ├── addon/                         # Addon interface + implementations
│   │   ├── addon.go
│   │   ├── mysql/
│   │   ├── postgres/
│   │   ├── redis/
│   │   └── meilisearch/
│   ├── compose/                       # Docker Compose types + YAML generation
│   ├── config/                        # Project config (envio.yaml)
│   ├── proxy/                         # Traefik proxy, TLS/mkcert integration
│   └── registry/                      # App/Addon discovery registry
```

---

## License

[MIT](LICENSE)
