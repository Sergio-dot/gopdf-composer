package main

import (
	"log"
	"os"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/engine"
)

func main() {
	log.Println("PDF Document Engine - Starting...")

	// 1. Load configuration via Viper
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Ensure output directory exists
	os.MkdirAll("output", 0o755)

	// 3. Initialize Engine with Config
	eng := engine.NewEngine(cfg)

	// 4. Generate PDF using default paths from config
	// Passing empty strings to use config defaults
	err = eng.GeneratePDF("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PDF Generated successfully at:", cfg.OutputPath)
}
