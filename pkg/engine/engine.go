package engine

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/Sergio-dot/gopdf-composer/config"
	"github.com/Sergio-dot/gopdf-composer/pkg/evaluator"
	"github.com/Sergio-dot/gopdf-composer/pkg/loader"
	"github.com/Sergio-dot/gopdf-composer/pkg/models"
	"github.com/Sergio-dot/gopdf-composer/pkg/renderer"
)

type Engine struct {
	config *config.Config
	loader loader.AssetLoader
}

func NewEngine(cfg *config.Config) *Engine {
	return &Engine{
		config: cfg,
		loader: loader.NewFileLoader(cfg.AssetDir),
	}
}

func (e *Engine) SetLoader(l loader.AssetLoader) {
	e.loader = l
}

func (e *Engine) GeneratePDF(controlFlowPath, runtimeContextPath, outputPath string) error {
	// If paths are empty, use those from config
	if controlFlowPath == "" {
		controlFlowPath = e.config.ControlFlowPath
	}
	if runtimeContextPath == "" {
		runtimeContextPath = e.config.RuntimeContextPath
	}
	if outputPath == "" {
		outputPath = e.config.OutputPath
	}

	// Load control flow
	cf, err := loader.LoadControlFlow(controlFlowPath)
	if err != nil {
		return err
	}

	// Load runtime context
	runtimeCtx, err := loader.LoadRuntimeContext(runtimeContextPath)
	if err != nil {
		return err
	}

	// Save PDF to file
	return e.GenerateToFile(cf, runtimeCtx, outputPath)
}

func (e *Engine) GenerateToFile(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext, outputPath string) error {
	renderer, err := e.render(cf, runtimeCtx)
	if err != nil {
		return err
	}
	return renderer.SaveToFile(outputPath)
}

func (e *Engine) GenerateToWriter(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext, w io.Writer) error {
	renderer, err := e.render(cf, runtimeCtx)
	if err != nil {
		return err
	}
	return renderer.WriteTo(w)
}

func (e *Engine) GenerateToBytes(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.GenerateToWriter(cf, runtimeCtx, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *Engine) render(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext) (*renderer.Renderer, error) {
	// Create renderer
	fontDir := e.config.FontDir
	if fontDir == "" {
		fontDir = filepath.Join(e.config.AssetDir, "fonts")
	}

	r := renderer.NewRenderer(runtimeCtx, fontDir, e.config.DefaultFont)
	pdf := r.GetPDF()

	// Load and register header assets
	headerBlocks, err := e.loadAssetBlocks(cf.Document.HeaderAssets)
	if err != nil {
		return nil, err
	}
	if len(headerBlocks) > 0 {
		pdf.SetHeaderFunc(func() {
			r.RenderHeader(headerBlocks)
		})
	}

	// Load and register footer assets
	footerBlocks, err := e.loadAssetBlocks(cf.Document.FooterAssets)
	if err != nil {
		return nil, err
	}
	if len(footerBlocks) > 0 {
		pdf.SetFooterFunc(func() {
			r.RenderFooter(footerBlocks, 15)
		})
	}

	// Process sections
	for _, section := range cf.Document.Structure {
		for _, assetRef := range section.Assets {
			// Evaluate condition
			shouldInclude, err := evaluator.Evaluate(assetRef.Conditions, runtimeCtx)
			if err != nil {
				return nil, err
			}

			if !shouldInclude {
				continue
			}

			// Load asset
			asset, err := e.loader.LoadAsset(assetRef.AssetID, assetRef.Version)
			if err != nil {
				return nil, err
			}

			// Render blocks
			for _, block := range asset.Blocks {
				if err := r.RenderBlock(&block); err != nil {
					return nil, err
				}
			}
		}
	}

	return r, nil
}

func (e *Engine) loadAssetBlocks(refs []models.AssetReference) ([]models.Block, error) {
	var allBlocks []models.Block
	for _, ref := range refs {
		asset, err := e.loader.LoadAsset(ref.AssetID, ref.Version)
		if err != nil {
			return nil, err
		}
		allBlocks = append(allBlocks, asset.Blocks...)
	}
	return allBlocks, nil
}
