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

	label := r.substituteVariables(props.Label)
	labelWidth := r.pdf.GetStringWidth(label)

	if props.Layout == "stacked" {
		labelGap := props.LabelLineGap
		if labelGap == 0 {
			labelGap = 8
		}

		sigWidth := props.SignatureWidth
		if sigWidth == 0 {
			sigWidth = available * 0.6
		}
		if sigWidth > available {
			sigWidth = available
		}

		lineStart := marginLeft
		switch props.Align {
		case "center":
			lineStart = marginLeft + (available-sigWidth)/2
		case "right":
			lineStart = marginLeft + available - sigWidth
		}

		y := r.pdf.GetY()
		if label != "" {
			labelX := lineStart + (sigWidth-labelWidth)/2
			if labelX < marginLeft {
				labelX = marginLeft
			}
			r.pdf.SetX(labelX)
			r.pdf.CellFormat(labelWidth, 5, label, "", 0, "L", false, 0, "")
			r.pdf.Ln(labelGap)
		}
		y = r.pdf.GetY()

		r.drawColor(props.LineColor)
		r.pdf.SetLineWidth(lineWidth)
		r.pdf.Line(lineStart, y, lineStart+sigWidth, y)
		r.drawColor("")
		r.pdf.SetY(y + 2)
	} else {
		// inline: label and signing line on the same baseline
		y := r.pdf.GetY()
		x := marginLeft
		if label != "" {
			r.pdf.SetX(x)
			r.pdf.CellFormat(0, 5, label, "", 0, "L", false, 0, "")
			x = marginLeft + labelWidth + 2
		}

		remaining := marginLeft + available - x
		sigWidth := props.SignatureWidth
		if sigWidth == 0 {
			sigWidth = remaining
		}
		if sigWidth > remaining {
			sigWidth = remaining
		}

		lineY := y + 4
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
