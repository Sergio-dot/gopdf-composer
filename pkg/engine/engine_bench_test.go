package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func BenchmarkGenerateToBytes(b *testing.B) {
	dir := b.TempDir()

	asset := models.Asset{
		Blocks: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Hello World", FontSize: 12}},
		},
	}
	assetData, _ := json.Marshal(asset)
	os.WriteFile(filepath.Join(dir, "greeting_v1.json"), assetData, 0644)

	rc := &models.RuntimeContext{Data: map[string]any{}}
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.GenerateToBytes(cf, rc)
	}
}
