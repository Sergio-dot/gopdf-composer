# Contributing

## Development workflow

1. Create a branch from `main` using the convention `<type>/<description>`:
   - `fix/` — bug fixes
   - `feat/` — new features
   - `refactor/` — structural changes without behavior changes
   - `test/` — test additions
   - `ci/` — CI/CD changes
   - `docs/` — documentation
   - `chore/` — maintenance tasks

2. Make your changes, keeping commits focused and well-described.

3. Run the full check suite locally before pushing:
   ```bash
   go vet ./...
   go test ./pkg/... ./config/ -race -coverprofile=coverage.out
   ```

4. Push your branch and open a pull request against `main`. CI will run `go vet`, tests with the race detector, and coverage checks on Go 1.24 and 1.25.

5. Pull requests require maintainer approval before merging. Once CI passes and your PR is reviewed and approved, it will be merged to `main`.

## Code style

- Match the existing code conventions in the file you're editing.
- No narrative comments (`// Load config`, `// Render table`). The code should speak for itself. Comments should explain *why*, not *what*.
- Run `go vet` and `staticcheck` before committing.

## Testing

- All packages under `pkg/` and `config/` should have test coverage.
- Use table-driven tests for exhaustive condition coverage.
- Use golden file hashes for PDF output regression tests.
- The `cmd/cli/` main package does not require unit tests (CLI entry point).

## Architecture

```
cmd/cli/        CLI entry point (slog logging, config loading)
config/         Viper-based configuration (YAML, .env, env vars)
pkg/
  models/       Data types (Asset, Block, ControlFlow, RuntimeContext)
  loader/       Asset loading with pluggable AssetLoader interface
  evaluator/    Conditional expression evaluation (==, !=, >, <, in, etc.)
  renderer/     PDF generation via gofpdf (one file per block type)
  engine/       Orchestrator for the full pipeline
```

## Project structure

- `config.yaml` — optional local overrides (gitignored)
- `.env` — optional environment overrides (gitignored)
- `examples/showcase/` — example flows and assets demonstrating all features
- `docs/adr/` — architecture decision records
