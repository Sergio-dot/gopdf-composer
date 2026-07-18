// Package renderer generates PDF documents from block structures using gofpdf.
// It supports text, image, table, container (row/column), pagebreak, loop,
// and line block types.
package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
	"github.com/phpdave11/gofpdf"
)

// Renderer holds the PDF document state and renders blocks into it.
type Renderer struct {
	pdf         *gofpdf.Fpdf
	context     *models.RuntimeContext
	defaultFont string
}

// Margins defines the page margins for a document.
type Margins struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// NewRenderer creates a Renderer with the given runtime context, font
// configuration, page orientation, page size, and optional custom margins.
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

// RenderBlock dispatches the block to its type-specific render function.
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

// RenderBlocks renders a sequence of blocks in order.
func (r *Renderer) RenderBlocks(blocks []models.Block) error {
	for _, block := range blocks {
		if err := r.RenderBlock(&block); err != nil {
			return err
		}
	}
	return nil
}

// RenderBlocksAtY renders blocks at a specific Y position, then restores the cursor.
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

// RenderHeader renders blocks as a page header.
func (r *Renderer) RenderHeader(blocks []models.Block) error {
	return r.RenderBlocks(blocks)
}

// RenderFooter renders blocks as a page footer at the given offset from the bottom.
func (r *Renderer) RenderFooter(blocks []models.Block, offsetFromBottom float64) error {
	r.pdf.SetY(-offsetFromBottom)
	return r.RenderBlocks(blocks)
}

// GetPDF returns the underlying gofpdf.Fpdf instance.
func (r *Renderer) GetPDF() *gofpdf.Fpdf {
	return r.pdf
}

// GetContext returns the RuntimeContext associated with this renderer.
func (r *Renderer) GetContext() *models.RuntimeContext {
	return r.context
}

// SetDefaultFont updates the default font used for rendering.
func (r *Renderer) SetDefaultFont(font string) {
	r.defaultFont = font
}
