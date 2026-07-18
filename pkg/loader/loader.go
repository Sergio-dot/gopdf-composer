// Package loader provides asset loading from filesystem and other sources.
// The AssetLoader interface enables custom loaders for S3, databases, etc.
package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

// AssetLoader is the interface for loading assets by ID and version.
// Implementations can load from filesystem, S3, database, or any source.
type AssetLoader interface {
	LoadAsset(assetID, version string) (*models.Asset, error)
}

// FileLoader loads assets from the local filesystem as {assetID}_v{version}.json.
type FileLoader struct {
	assetDir string
}

// NewFileLoader creates a FileLoader that reads assets from the given directory.
func NewFileLoader(assetDir string) *FileLoader {
	return &FileLoader{assetDir: assetDir}
}

func (l *FileLoader) LoadAsset(assetID, version string) (*models.Asset, error) {
	filename := fmt.Sprintf("%s_v%s.json", assetID, version)
	path := filepath.Join(l.assetDir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var asset models.Asset
	err = json.Unmarshal(data, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// LoadControlFlow reads and parses a ControlFlow JSON file from the given path.
func LoadControlFlow(path string) (*models.ControlFlow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cf models.ControlFlow
	err = json.Unmarshal(data, &cf)
	if err != nil {
		return nil, err
	}

	return &cf, nil
}

// LoadRuntimeContext reads and parses a runtime context JSON file from the given path.
func LoadRuntimeContext(path string) (*models.RuntimeContext, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var contextData map[string]any
	err = json.Unmarshal(data, &contextData)
	if err != nil {
		return nil, err
	}

	return &models.RuntimeContext{Data: contextData}, nil
}
