// Package engine orchestrates the PDF generation pipeline: loading assets,
// evaluating conditions, rendering blocks, and producing output.
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

// Engine is the primary orchestrator for PDF generation. It holds configuration
// and an AssetLoader, and provides methods to generate PDFs to files, writers,
// or byte slices.
type Engine struct {
	config *config.Config
	loader loader.AssetLoader
}

// NewEngine creates an Engine with the given config and a filesystem-based asset loader.
func NewEngine(cfg *config.Config) *Engine {
	return &Engine{
		config: cfg,
		loader: loader.NewFileLoader(cfg.AssetDir),
	}
}

// SetLoader replaces the default filesystem asset loader with a custom implementation.
func (e *Engine) SetLoader(l loader.AssetLoader) {
	e.loader = l
}

// GeneratePDF loads control flow and runtime context from files, generates the
// PDF, and writes it to the output path. Empty path arguments fall back to config defaults.
func (e *Engine) GeneratePDF(controlFlowPath, runtimeContextPath, outputPath string) error {
	if controlFlowPath == "" {
		controlFlowPath = e.config.ControlFlowPath
	}
	if runtimeContextPath == "" {
		runtimeContextPath = e.config.RuntimeContextPath
	}
	if outputPath == "" {
		outputPath = e.config.OutputPath
	}

	cf, err := loader.LoadControlFlow(controlFlowPath)
	if err != nil {
		return err
	}

	runtimeCtx, err := loader.LoadRuntimeContext(runtimeContextPath)
	if err != nil {
		return err
	}

	return e.GenerateToFile(cf, runtimeCtx, outputPath)
}

// GenerateToFile renders the given control flow and runtime context and saves
// the resulting PDF to the specified path.
func (e *Engine) GenerateToFile(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext, outputPath string) error {
	renderer, err := e.render(cf, runtimeCtx)
	if err != nil {
		return err
	}
	return renderer.SaveToFile(outputPath)
}

// GenerateToWriter renders the given control flow and runtime context and writes
// the resulting PDF to the provided io.Writer.
func (e *Engine) GenerateToWriter(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext, w io.Writer) error {
	renderer, err := e.render(cf, runtimeCtx)
	if err != nil {
		return err
	}
	_, err = renderer.WriteTo(w)
	return err
}

// GenerateToBytes renders the given control flow and runtime context and returns
// the resulting PDF as a byte slice. Suitable for HTTP responses.
func (e *Engine) GenerateToBytes(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.GenerateToWriter(cf, runtimeCtx, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *Engine) render(cf *models.ControlFlow, runtimeCtx *models.RuntimeContext) (*renderer.Renderer, error) {
	fontDir := e.config.FontDir
	if fontDir == "" && e.config.AssetDir != "" {
		fontDir = filepath.Join(e.config.AssetDir, "fonts")
	}

	r := renderer.NewRenderer(runtimeCtx, fontDir, e.config.DefaultFont,
		cf.Document.Orientation, cf.Document.PageSize,
		toMargins(&cf.Document),
	)
	pdf := r.GetPDF()

	for family, styles := range e.config.FontFiles {
		for style, path := range styles {
			pdf.AddUTF8Font(family, style, path)
			if pdf.Error() != nil {
				pdf.ClearError()
			}
		}
	}
	if len(e.config.FontFiles) > 0 && e.config.DefaultFont == "" {
		for name := range e.config.FontFiles {
			r.SetDefaultFont(name)
			break
		}
	}

	headerBlocks, err := e.loadAssetBlocks(cf.Document.HeaderAssets)
	if err != nil {
		return nil, err
	}
	if len(headerBlocks) > 0 {
		pdf.SetHeaderFunc(func() {
			r.RenderHeader(headerBlocks)
		})
	}

	footerBlocks, err := e.loadAssetBlocks(cf.Document.FooterAssets)
	if err != nil {
		return nil, err
	}
	if len(footerBlocks) > 0 {
		pdf.SetFooterFunc(func() {
			r.RenderFooter(footerBlocks, 15)
		})
	}

	for _, section := range cf.Document.Structure {
		for _, assetRef := range section.Assets {
			shouldInclude, err := evaluator.Evaluate(assetRef.Conditions, runtimeCtx)
			if err != nil {
				return nil, err
			}

			if !shouldInclude {
				continue
			}

			asset, err := e.loader.LoadAsset(assetRef.AssetID, assetRef.Version)
			if err != nil {
				return nil, err
			}

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

func toMargins(doc *models.Document) *renderer.Margins {
	if doc.MarginLeft == 0 && doc.MarginTop == 0 && doc.MarginRight == 0 && doc.MarginBottom == 0 {
		return nil
	}
	return &renderer.Margins{
		Left:   doc.MarginLeft,
		Top:    doc.MarginTop,
		Right:  doc.MarginRight,
		Bottom: doc.MarginBottom,
	}
}
