# Envio

A CLI tool for creating and managing local Docker development environments.

Envio generates Docker Compose setups for various application types with optional addon services. It's built to be extensible — adding new apps or addons is as simple as implementing an interface.

## Install

```bash
go install github.com/danjdewhurst/envio@latest
```

Or build from source:

```bash
git clone https://github.com/danjdewhurst/envio.git
cd envio
go build -o envio .
```

## Quick Start

```bash
# Initialize a Laravel project with MySQL and Redis
envio init laravel --addon mysql --addon redis

# Start the environment
envio up

# Check status
envio status

# Stop the environment
envio down
```

## Commands

| Command | Description |
|---------|-------------|
| `envio init <app>` | Initialize a new Docker environment |
| `envio up` | Start the Docker environment (`docker compose up -d`) |
| `envio down` | Stop the Docker environment (`docker compose down`) |
| `envio status` | Show app config and running container status |
| `envio apps` | List available application types |
| `envio addons` | List available addons |

### `envio init`

```bash
envio init <app> [--addon <name>]...
```

Generates a `docker-compose.yml` and `envio.yaml` in the current directory.

**Flags:**
- `--addon, -a` — Add one or more addons (repeatable)

## Supported Apps

| App | Description |
|-----|-------------|
| `laravel` | PHP Laravel with Nginx and PHP-FPM |

## Supported Addons

| Addon | Description |
|-------|-------------|
| `mysql` | MySQL 8.0 |
| `postgres` | PostgreSQL 16 |
| `redis` | Redis |
| `meilisearch` | Meilisearch |

## Extending Envio

### Adding a new app

1. Create a new package under `internal/app/` (e.g. `internal/app/wordpress/`)
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

1. Create a new package under `internal/addon/` (e.g. `internal/addon/mailpit/`)
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
│   └── registry/                      # App/Addon discovery registry
```

## License

See [LICENSE](LICENSE).
