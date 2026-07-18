package main

import (
	"log/slog"
	"os"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/engine"
)

func main() {
	slog.Info("PDF Document Engine - Starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	os.MkdirAll("output", 0o755)

	eng := engine.NewEngine(cfg)

	err = eng.GeneratePDF("", "", "")
	if err != nil {
		slog.Error("PDF generation failed", "error", err)
		os.Exit(1)
	}

	slog.Info("PDF generated successfully", "output", cfg.OutputPath)
}
