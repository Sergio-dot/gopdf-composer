package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func (r *Renderer) renderSignature(block *models.Block) error {
	props := block.SignatureProperties
	if props == nil {
		return fmt.Errorf("signature block missing signatureProperties")
	}

	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	if props.Label != "" {
		r.textColor("")
		r.pdf.SetFont(r.defaultFont, "", 9)
		r.pdf.CellFormat(0, 5, r.substituteVariables(props.Label), "", 0, "L", false, 0, "")
		r.pdf.Ln(5)
	}

	pageWidth, _ := r.pdf.GetPageSize()
	marginLeft, _, marginRight, _ := r.pdf.GetMargins()
	width := pageWidth - marginLeft - marginRight
	if props.SignatureWidth > 0 {
		width = props.SignatureWidth
	}

	lineWidth := props.LineWidth
	if lineWidth <= 0 {
		lineWidth = 0.4
	}

	y := r.pdf.GetY()
	r.drawColor(props.LineColor)
	r.pdf.SetLineWidth(lineWidth)
	r.pdf.Line(marginLeft, y, marginLeft+width, y)
	r.drawColor("")
	r.pdf.SetY(y + 2)

	if props.MarginBottom > 0 {
		r.pdf.Ln(props.MarginBottom)
	}

	return nil
}
