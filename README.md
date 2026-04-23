# gopdf-composer

A flexible, agnostic PDF generation engine for Go. It allows you to generate complex documents by defining control flows and modular assets in JSON, supporting variable substitution and conditional rendering.

## Features

- **JSON-Driven**: Define your document structure and content blocks in JSON.
- **Conditional Rendering**: Use logical conditions to include or exclude sections of the document based on runtime data.
- **Variable Substitution**: Inject data into your documents using `{{variable}}` syntax.
- **Modular Assets**: Manage reusable content blocks (Text, Images, Tables, Containers).
- **Flexible Configuration**: Configure via environment variables, YAML, or `.env` files using Viper.
- **Library Friendly**: Can be used as a CLI tool or imported as a Go module for APIs.

## Installation

```bash
go get github.com/Sergio-dot/gopdf-composer
```

## Configuration

The engine uses **Viper** for configuration. It looks for settings in the following order:

1.  **Environment Variables**: Prefixed with `GOPDF_` (e.g., `GOPDF_ASSET_DIR`).
2.  **`.env` File**: Key-value pairs (e.g., `ASSET_DIR=assets`).
3.  **`config.yaml`**: Standard YAML format.

### Available Settings

| Key | Description | Default |
|-----|-------------|---------|
| `asset_dir` | Directory containing JSON assets | `assets/` |
| `control_flow_path` | Path to the document structure JSON | `flows/section_oriented_control_flow.json` |
| `runtime_context_path` | Path to the runtime data JSON | `contexts/runtime_context.json` |
| `output_path` | Where to save the generated PDF | `output/document.pdf` |
| `font_dir` | Directory containing TTF fonts | `assets/fonts` |

## Usage

### As a CLI Tool

1. Clone the repo and build:
   ```bash
   go build -o gopdf main.go
   ```
2. Run with default config:
   ```bash
   ./gopdf
   ```
3. Run with custom environment:
   ```bash
   GOPDF_OUTPUT_PATH="my_report.pdf" ./gopdf
   ```

### As a Library (API)

```go
import (
    "github.com/Sergio-dot/gopdf-composer/config"
    "github.com/Sergio-dot/gopdf-composer/pkg/engine"
    "github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Load basic config
    cfg, _ := config.LoadConfig()
    eng := engine.NewEngine(cfg)

    // 2. Prepare data (usually from DB or Request)
    flow := &models.ControlFlow{ /* ... */ }
    ctx := &models.RuntimeContext{Data: map[string]any{"user": "Sergio"}}

    // 3. Generate to bytes for HTTP response
    pdfBytes, err := eng.GenerateToBytes(flow, ctx)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Write(pdfBytes)
}
```

## Custom Asset Loaders

You can implement the `AssetLoader` interface to load assets from S3, a Database, or any other source:

```go
type MyCustomLoader struct {}

func (l *MyCustomLoader) LoadAsset(id, version string) (*models.Asset, error) {
    // Fetch from Database...
}

// ... in your code
eng.SetLoader(&MyCustomLoader{})
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
