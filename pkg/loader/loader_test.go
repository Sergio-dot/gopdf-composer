package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func TestFileLoaderLoadAsset(t *testing.T) {
	dir := t.TempDir()

	asset := models.Asset{
		Blocks: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Hello", FontSize: 12}},
		},
	}
	data, err := json.Marshal(asset)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "header_v1.json"), data, 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewFileLoader(dir)
	loaded, err := loader.LoadAsset("header", "1")
	if err != nil {
		t.Fatal(err)
	}

	if len(loaded.Blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(loaded.Blocks))
	}
	if loaded.Blocks[0].Type != "text" {
		t.Errorf("expected text block, got %s", loaded.Blocks[0].Type)
	}
}

func TestFileLoaderLoadAssetNotFound(t *testing.T) {
	dir := t.TempDir()
	loader := NewFileLoader(dir)
	_, err := loader.LoadAsset("missing", "1")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadControlFlow(t *testing.T) {
	dir := t.TempDir()

	cf := models.ControlFlow{
		Document: models.Document{
			Structure: []models.Section{
				{Assets: []models.AssetReference{
					{AssetID: "body", Version: "1"},
				}},
			},
		},
	}
	data, err := json.Marshal(cf)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "flow.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadControlFlow(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Document.Structure) != 1 {
		t.Errorf("expected 1 section, got %d", len(loaded.Document.Structure))
	}
}

func TestLoadControlFlowInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadControlFlow(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadRuntimeContext(t *testing.T) {
	dir := t.TempDir()

	ctxData := map[string]any{
		"name": "Sergio",
		"age":  30,
	}
	data, err := json.Marshal(ctxData)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "context.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadRuntimeContext(path)
	if err != nil {
		t.Fatal(err)
	}

	if val, ok := loaded.Get("name"); !ok || val != "Sergio" {
		t.Errorf("expected 'Sergio', got %v", val)
	}
}

func TestLoadRuntimeContextInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadRuntimeContext(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
