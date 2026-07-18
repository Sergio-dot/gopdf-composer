package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
	"github.com/phpdave11/gofpdf"
)

type Renderer struct {
	pdf         *gofpdf.Fpdf
	context     *models.RuntimeContext
	defaultFont string
}

type Margins struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

func NewRenderer(runtimeCtx *models.RuntimeContext, fontDir, defaultFont, orientation, pageSize string, margins *Margins) *Renderer {
	if orientation == "" {
		orientation = "P"
	}
	if pageSize == "" {
		pageSize = "A4"
	}
	if defaultFont == "" {
		defaultFont = "Arial"
	}

	pdf := gofpdf.New(orientation, "mm", pageSize, fontDir)

	if margins != nil {
		pdf.SetMargins(margins.Left, margins.Top, margins.Right)
		if margins.Bottom != 0 {
			pdf.SetAutoPageBreak(true, margins.Bottom)
		}
	}

	pdf.SetFont(defaultFont, "", 12)
	pdf.AddPage()
	pdf.AliasNbPages("{nb}")

	return &Renderer{pdf: pdf, context: runtimeCtx, defaultFont: defaultFont}
}

func (r *Renderer) RenderBlock(block *models.Block) error {
	switch block.Type {
	case "text":
		return r.renderText(block)
	case "image":
		return r.renderImage(block)
	case "table":
		return r.renderTable(block)
	case "container":
		return r.renderContainer(block)
	case "pagebreak":
		return r.renderPageBreak(block)
	case "loop":
		return r.renderLoop(block)
	case "line":
		return r.renderLine(block)
	default:
		return fmt.Errorf("unknown block type: %s", block.Type)
	}
}

func (r *Renderer) RenderBlocks(blocks []models.Block) error {
	for _, block := range blocks {
		if err := r.RenderBlock(&block); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) RenderBlocksAtY(y float64, blocks []models.Block) error {
	origX := r.pdf.GetX()
	origY := r.pdf.GetY()

	r.pdf.SetXY(origX, y)
	if err := r.RenderBlocks(blocks); err != nil {
		r.pdf.SetXY(origX, origY)
		return err
	}

	r.pdf.SetXY(origX, origY)
	return nil
}

func (r *Renderer) RenderHeader(blocks []models.Block) error {
	return r.RenderBlocks(blocks)
}

func (r *Renderer) RenderFooter(blocks []models.Block, offsetFromBottom float64) error {
	r.pdf.SetY(-offsetFromBottom)
	return r.RenderBlocks(blocks)
}

func (r *Renderer) GetPDF() *gofpdf.Fpdf {
	return r.pdf
}

func (r *Renderer) GetContext() *models.RuntimeContext {
	return r.context
}
