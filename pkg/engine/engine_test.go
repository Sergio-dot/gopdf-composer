package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/loader"
	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func TestEngineGenerateToBytes(t *testing.T) {
	dir := t.TempDir()

	// Create a minimal asset
	asset := models.Asset{
		Blocks: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Hello World", FontSize: 12}},
		},
	}
	assetData, _ := json.Marshal(asset)
	os.WriteFile(filepath.Join(dir, "greeting_v1.json"), assetData, 0644)

	// Create runtime context
	rc := &models.RuntimeContext{Data: map[string]any{}}

	// Create control flow pointing to the asset
	cf := &models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{
					{AssetID: "greeting", Version: "1"},
				}},
			},
		},
	}

	cfg := &config.Config{
		AssetDir:    dir,
		FontDir:     dir,
		DefaultFont: "Arial",
	}

	eng := NewEngine(cfg)

	pdfBytes, err := eng.GenerateToBytes(cf, rc)
	if err != nil {
		t.Fatalf("GenerateToBytes failed: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output")
	}

	// PDF magic bytes
	magic := "%PDF-"
	if string(pdfBytes[:len(magic)]) != magic {
		t.Errorf("PDF magic bytes not found, got: %s", string(pdfBytes[:min(len(magic), len(pdfBytes))]))
	}
}

func TestEngineGenerateToBytesEmpty(t *testing.T) {
	dir := t.TempDir()

	cf := &models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{}},
			},
		},
	}
	rc := &models.RuntimeContext{Data: map[string]any{}}

	cfg := &config.Config{
		AssetDir:    dir,
		FontDir:     dir,
		DefaultFont: "Arial",
	}

	eng := NewEngine(cfg)
	pdfBytes, err := eng.GenerateToBytes(cf, rc)
	if err != nil {
		t.Fatalf("GenerateToBytes failed: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output even with no content")
	}
}

func TestEngineSetLoader(t *testing.T) {
	dir := t.TempDir()

	// Create asset with a custom loader
	asset := models.Asset{
		Blocks: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Custom Loaded", FontSize: 12}},
		},
	}

	eng := NewEngine(&config.Config{
		AssetDir:    dir,
		FontDir:     dir,
		DefaultFont: "Arial",
	})

	eng.SetLoader(&stubLoader{asset: &asset})

	cf := &models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{
					{AssetID: "any", Version: "1"},
				}},
			},
		},
	}
	rc := &models.RuntimeContext{Data: map[string]any{}}

	pdfBytes, err := eng.GenerateToBytes(cf, rc)
	if err != nil {
		t.Fatalf("GenerateToBytes with custom loader failed: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
}

func TestEngineGenerateToBytesMissingAsset(t *testing.T) {
	dir := t.TempDir()

	cf := &models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{
					{AssetID: "missing", Version: "1"},
				}},
			},
		},
	}
	rc := &models.RuntimeContext{Data: map[string]any{}}

	cfg := &config.Config{
		AssetDir:    dir,
		FontDir:     dir,
		DefaultFont: "Arial",
	}

	eng := NewEngine(cfg)
	_, err := eng.GenerateToBytes(cf, rc)
	if err == nil {
		t.Error("expected error for missing asset")
	}
}

func TestEngineGenerateToBytesConditionalSkip(t *testing.T) {
	dir := t.TempDir()

	asset := models.Asset{
		Blocks: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Should not appear", FontSize: 12}},
		},
	}
	assetData, _ := json.Marshal(asset)
	os.WriteFile(filepath.Join(dir, "conditional_v1.json"), assetData, 0644)

	cf := &models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{
					{
						AssetID: "conditional",
						Version: "1",
						Conditions: &models.Condition{
							Field: "show",
							Op:    "==",
							Value: false,
						},
					},
				}},
			},
		},
	}
	rc := &models.RuntimeContext{Data: map[string]any{"show": true}}

	cfg := &config.Config{
		AssetDir:    dir,
		FontDir:     dir,
		DefaultFont: "Arial",
	}

	eng := NewEngine(cfg)
	pdfBytes, err := eng.GenerateToBytes(cf, rc)
	if err != nil {
		t.Fatalf("GenerateToBytes failed: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("expected non-empty PDF output")
	}
}

type stubLoader struct {
	asset *models.Asset
}

func (s *stubLoader) LoadAsset(assetID, version string) (*models.Asset, error) {
	return s.asset, nil
}

var _ loader.AssetLoader = (*stubLoader)(nil)
