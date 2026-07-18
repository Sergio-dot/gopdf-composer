package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfigDefaults(t *testing.T) {
	viper.Reset()

	// Viper is global state, so isolate with a temp working directory
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AssetDir != "assets/" {
		t.Errorf("expected default asset_dir 'assets/', got %q", cfg.AssetDir)
	}
	if cfg.ControlFlowPath != "flows/flow.json" {
		t.Errorf("expected default control_flow_path 'flows/flow.json', got %q", cfg.ControlFlowPath)
	}
	if cfg.OutputPath != "output/document.pdf" {
		t.Errorf("expected default output_path 'output/document.pdf', got %q", cfg.OutputPath)
	}
	if cfg.DefaultFont != "Arial" {
		t.Errorf("expected default font 'Arial', got %q", cfg.DefaultFont)
	}
}

func TestLoadConfigFromYAML(t *testing.T) {
	viper.Reset()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	content := `
asset_dir: custom_assets/
control_flow_path: custom_flow.json
output_path: custom_output.pdf
default_font: Helvetica
`
	os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AssetDir != "custom_assets/" {
		t.Errorf("expected 'custom_assets/', got %q", cfg.AssetDir)
	}
	if cfg.ControlFlowPath != "custom_flow.json" {
		t.Errorf("expected 'custom_flow.json', got %q", cfg.ControlFlowPath)
	}
	if cfg.DefaultFont != "Helvetica" {
		t.Errorf("expected 'Helvetica', got %q", cfg.DefaultFont)
	}
}

func TestLoadConfigFromEnvFile(t *testing.T) {
	viper.Reset()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	content := `
ASSET_DIR=env_assets/
`
	os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AssetDir != "env_assets/" {
		t.Errorf("expected 'env_assets/', got %q", cfg.AssetDir)
	}
}

func TestLoadConfigFromEnvVar(t *testing.T) {
	viper.Reset()

	dir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	os.Setenv("GOPDF_ASSET_DIR", "env_var_assets/")
	defer os.Unsetenv("GOPDF_ASSET_DIR")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AssetDir != "env_var_assets/" {
		t.Errorf("expected 'env_var_assets/', got %q", cfg.AssetDir)
	}
}
