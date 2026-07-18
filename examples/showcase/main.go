package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/engine"
	"github.com/Sergio-dot/gopdf-composer/pkg/loader"
)

func main() {
	flowsDir := "flows"
	contextsDir := "contexts"
	outputDir := "output"
	assetsDir := "assets"
	fontsDir := "assets/fonts"

	os.MkdirAll(outputDir, 0755)

	flowFiles := []struct {
		name  string
		label string
	}{
		{"a3-showcase", "A3 • Portrait"},
		{"a4-showcase", "A4 • Portrait"},
		{"a5-showcase", "A5 • Landscape"},
		{"letter-showcase", "Letter • Portrait"},
		{"legal-showcase", "Legal • Portrait"},
	}

	ctx, err := loader.LoadRuntimeContext(filepath.Join(contextsDir, "showcase-context.json"))
	if err != nil {
		slog.Error("failed to load context", "error", err)
		os.Exit(1)
	}

	cfg := &config.Config{
		AssetDir:    assetsDir,
		FontDir:     fontsDir,
		DefaultFont: "Arial",
	}

	for _, ff := range flowFiles {
		cf, err := loader.LoadControlFlow(filepath.Join(flowsDir, ff.name+".json"))
		if err != nil {
			slog.Error("failed to load control flow", "flow", ff.name, "error", err)
			os.Exit(1)
		}

		ctx.Set("sizeLabel", ff.label)

		eng := engine.NewEngine(cfg)

		outPath := filepath.Join(outputDir, ff.name+".pdf")
		if err := eng.GenerateToFile(cf, ctx, outPath); err != nil {
			slog.Error("failed to generate PDF", "flow", ff.name, "error", err)
			os.Exit(1)
		}
		slog.Info("generated", "file", outPath, "size", ff.label)
	}
}
