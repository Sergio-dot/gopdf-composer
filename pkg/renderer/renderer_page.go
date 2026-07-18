package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func (r *Renderer) renderPageBreak(_ *models.Block) error {
	r.pdf.AddPage()
	return nil
}

func (r *Renderer) renderLine(block *models.Block) error {
	props := block.LineProperties
	if props == nil {
		return fmt.Errorf("line block missing lineProperties")
	}

	if props.Margin > 0 {
		r.pdf.Ln(props.Margin)
	}

	pageWidth, _ := r.pdf.GetPageSize()
	marginLeft, _, marginRight, _ := r.pdf.GetMargins()
	availableWidth := pageWidth - marginLeft - marginRight

	lineWidth := props.Width
	if lineWidth <= 0 {
		lineWidth = 0.4
	}

	y := r.pdf.GetY()

	r.drawColor(props.Color)
	r.pdf.SetLineWidth(lineWidth)
	r.pdf.Line(marginLeft, y, marginLeft+availableWidth, y)
	r.drawColor("")

	if props.Margin > 0 {
		r.pdf.Ln(props.Margin)
	}

	return nil
}
