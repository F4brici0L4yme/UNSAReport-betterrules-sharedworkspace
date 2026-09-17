# AGENTS.md - UNSAReport CLI

Go CLI that scaffolds Typst-based lab reports, manages versioned templates/components, captures terminal output, and compiles submissions.

## Commands

```bash
# Build
go build ./cmd/unsarep

# Unit tests (integration tests are skipped by default)
go test ./...

# Integration tests — need external tools (typst, freeze, magick) on PATH
go test -tags integration ./...

# Single test / package
go test -run TestName ./internal/services/...

# Lint (gofmt + goimports enforced; revive/misspell)
golangci-lint run
```

## Architecture

Ports and Adapters with dependency injection.

* `internal/ports/` — interfaces (Compiler, Archiver, Fetcher, Renderer, etc.)
* `internal/adapters/` — concrete implementations (osfs, github, config, registry, typst, freeze, zipper)
* `internal/services/` — business logic (install, update, prepare, capture, component, share)
* `internal/cmd/` — cobra commands
* `internal/mocks/` — testify mocks for port interfaces
* `internal/dependencies/` — runtime checks for external tools

Entry point: `cmd/unsarep/main.go` → `internal/cmd.Execute()`.

## Key Facts

* **Version control is git, not jj.** Repo is a fork: `origin` → `F4brici0L4yme/UNSAReport-betterrules-sharedworkspace`, `upstream` → `UNSAReport/UNSAReport`.
* **Version injection** via ldflags: `-X github.com/UNSAReport/UNSAReport/internal/ports.Version={{.Version}}` (this is what `flake.nix` uses; `internal/cmd.Version` is just an alias of `ports.Version`). Default version lives in `internal/ports/version.go`.
* **External tools** validated at runtime (`internal/dependencies`): Typst, Freeze (charmbracelet), ImageMagick — note the ImageMagick binary is `magick`, not `convert`.
* **Goroutine leak detection**: every test package has a `main_test.go` calling `goleak.VerifyTestMain`.
* **Integration tests** carry `//go:build integration` and live in several packages (services, adapters/zipper, osfs, typst, freeze, github), not just `internal/services`.
* **Config**: viper with `UNSAREP_` env prefix; project config is `unsareport.json` plus an auto-maintained `unsareport.lock`.
* **Nix dev shell**: `nix develop` / `direnv allow`. The shell sets a project-local `GOPATH=$PWD/.go` (gitignored as `.go`).
* **`unsarep share`** creates a disposable GitHub repo from the current dir via the `gh` CLI (git-init + initial commit when needed, private by default).
* **Agent skills** in `.agents/skills/`: `report-header`, `report-content`, `report-code`, `report-terminal`, `report-references`, `report-review`, `browser-capture`, `puml-diagram`. The report header/tables/content live in the Typst template (separate `UNSAReport/templates` repo), not in this Go CLI. See `docs/guia-estudiante.md` for a Spanish quick-start guide.

## Testing

* Unit tests: `testify` (assert/require) + `mock` for ports; `t.Parallel()` for independent tests, `t.TempDir()` for filesystem cleanup.
* Integration tests additionally require network for the `adapters/github` tests (GitHub API) and the external tools on PATH.

## Build & Release

```bash
# Development build
go build ./cmd/unsarep

# Release build with version (matches flake.nix ldflags)
go build -ldflags "-X github.com/UNSAReport/UNSAReport/internal/ports.Version=1.0.0" ./cmd/unsarep

# Nix build
nix build
```

### Updating the Nix Vendor Hash

When Go dependencies change, the Nix build fails on a mismatched `vendorHash`:

1. Set `vendorHash = "";` (or `lib.fakeHash`) in `flake.nix`.
2. Run `nix build -L`.
3. Copy the `got:` hash from the failure output back into `vendorHash`.

## Conventions

* Error handling: `samber/oops` for stack traces in debug mode.
* CLI framework: cobra + viper config binding.
* Logging: `log/slog` text handler to stderr.
* Linters: revive (exported/package-comments disabled), misspell.
