package main

import (
	"log"
	"os"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/engine"
)

func main() {
	log.Println("PDF Document Engine - Starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	os.MkdirAll("output", 0o755)

	eng := engine.NewEngine(cfg)

	err = eng.GeneratePDF("", "", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PDF Generated successfully at:", cfg.OutputPath)
}
