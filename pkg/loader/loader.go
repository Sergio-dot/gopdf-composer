package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)
type AssetLoader interface {
	LoadAsset(assetID, version string) (*models.Asset, error)
}

type FileLoader struct {
	assetDir string
}

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
