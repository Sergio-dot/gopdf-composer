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

	fontFamily := r.defaultFont
	if props.FontFamily != "" {
		fontFamily = props.FontFamily
	}

	fontStyle := ""
	switch props.FontWeight {
	case "bold":
		fontStyle = "B"
	case "italic":
		fontStyle = "I"
	case "boldItalic", "bold_italic":
		fontStyle = "BI"
	}

	fontSize := props.FontSize
	if fontSize == 0 {
		fontSize = 9
	}

	r.textColor(props.FontColor)
	r.pdf.SetFont(fontFamily, fontStyle, fontSize)

	pageWidth, _ := r.pdf.GetPageSize()
	marginLeft, _, marginRight, _ := r.pdf.GetMargins()
	available := pageWidth - marginLeft - marginRight

	lineWidth := props.LineWidth
	if lineWidth <= 0 {
		lineWidth = 0.4
	}

	sigWidth := available
	if props.SignatureWidth > 0 {
		sigWidth = props.SignatureWidth
	}

	align := "L"
	switch props.Align {
	case "right":
		align = "R"
	case "center":
		align = "C"
	}

	label := r.substituteVariables(props.Label)

	if props.Layout == "below" {
		if label != "" {
			r.pdf.CellFormat(0, 5, label, "", 0, align, false, 0, "")
			r.pdf.Ln(7)
		}
		y := r.pdf.GetY()
		r.drawColor(props.LineColor)
		r.pdf.SetLineWidth(lineWidth)
		r.pdf.Line(marginLeft, y, marginLeft+sigWidth, y)
		r.drawColor("")
		r.pdf.SetY(y + 2)
	} else {
		// inline: label and signing line on the same baseline
		y := r.pdf.GetY()
		x := marginLeft
		if label != "" {
			r.pdf.SetX(x)
			r.pdf.CellFormat(0, 5, label, "", 0, "L", false, 0, "")
			x += r.pdf.GetStringWidth(label) + 2
		}
		lineY := y + 2.5
		r.drawColor(props.LineColor)
		r.pdf.SetLineWidth(lineWidth)
		r.pdf.Line(x, lineY, x+sigWidth, lineY)
		r.drawColor("")
		r.pdf.SetY(y + 5)
	}

	if props.MarginBottom > 0 {
		r.pdf.Ln(props.MarginBottom)
	}

	return nil
}
