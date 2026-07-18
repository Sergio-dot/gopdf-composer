# Changelog

## [Unreleased]

### Added
- Compound condition logic (`and`, `or`, `not`) in the evaluator.
- Numeric comparison operators (`>`, `<`, `>=`, `<=`).
- String and array operators (`contains`, `in`).
- Configurable page size, orientation, and margins via Document model.
- Horizontal line block type with color, width, and margin.
- Example showcase with 5 page sizes (A3, A4, A5, Letter, Legal).
- Architecture decision record for PDF library choice.
- Benchmark tests for renderer and engine hot paths.

### Changed
- Replaced `log` package with `log/slog` in CLI entry point.
- Split `renderer.go` (634 lines) into 9 focused files by block type.

### Fixed
- Nil pointer dereference when `HeaderStyle` is nil in table rendering.
- Pre-compiled variable substitution regex to package level.
- `.env` file parsing was silently failing due to missing `SetConfigType("env")`.
- Go 1.24 coverage tooling error by excluding `cmd/cli/` from coverage target.

## [0.1.0] — Initial

### Added
- JSON-driven control flow for document structure definition.
- Modular asset system with text, image, table, container, pagebreak, and loop blocks.
- Runtime context with dot-notation nested variable access.
- Conditional rendering with `==` and `!=` operators.
- Variable substitution with `{{variable}}` syntax.
- Page numbering (`{{page}}`, `{{totalPages}}`).
- Per-page header and footer support via control flow assets.
- Dynamic table rows via `RowsDataSource` from context arrays.
- Pluggable `AssetLoader` interface for custom asset sources (S3, DB, etc.).
- 12-factor configuration via Viper (YAML, `.env`, environment variables).
- Three output modes: file, `io.Writer`, `[]byte` (HTTP-ready).
- MIT License.
