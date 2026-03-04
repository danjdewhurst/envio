# CLAUDE.md

## Project Overview

Envio is a Go CLI tool that generates and manages local Docker development environments. It uses a plugin-like architecture where apps (Laravel, WordPress, etc.) and addons (Redis, MySQL, etc.) are registered via interfaces.

## Build & Run

```bash
# Build
go build -o envio .

# Run
go run .

# Build check (all packages)
go build ./...
```

## Architecture

- **Go + Cobra** for CLI framework
- **Interface-based extensibility**: `app.App` and `addon.Addon` interfaces in `internal/app/app.go` and `internal/addon/addon.go`
- **Registry pattern**: `internal/registry/` discovers and manages all apps/addons. New implementations are registered in `internal/registry/defaults.go`
- **Compose generation**: `internal/compose/` handles Docker Compose YAML file creation from service definitions
- **Project config**: `internal/config/` manages `envio.yaml` which tracks the app type and enabled addons per project

## Key Conventions

- Each app lives in its own package: `internal/app/<name>/`
- Each addon lives in its own package: `internal/addon/<name>/`
- CLI commands are in `cmd/` — one file per command
- All new apps/addons must be registered in `internal/registry/defaults.go`
- Services, volumes, and networks use the types defined in `internal/compose/types.go`
- The `envio` network (bridge driver) is added to all compose setups by the `init` command

## Adding New Apps/Addons

1. Create a new package implementing the `App` or `Addon` interface
2. Register it in `internal/registry/defaults.go`
3. For apps: list compatible addon names in `AvailableAddons()`
