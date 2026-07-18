package renderer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

var varRe = regexp.MustCompile(`\{\{([\w.]+)\}\}`)

func (r *Renderer) renderText(block *models.Block) error {
	if block.TextProperties == nil {
		return fmt.Errorf("text block missing textProperties")
	}

	props := block.TextProperties

	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	fontFamily := props.FontFamily
	if fontFamily == "" {
		fontFamily = r.defaultFont
	}

	fontWeight := ""
	switch props.FontWeight {
	case "bold":
		fontWeight = "B"
	case "italic":
		fontWeight = "I"
	}
	r.pdf.SetFont(fontFamily, fontWeight, props.FontSize)

	r.textColor(props.FontColor)

	text := r.substituteVariables(props.Text)

	align := "L"
	switch props.Align {
	case "center":
		align = "C"
	case "right":
		align = "R"
	}

	if props.BackgroundColor != "" {
		pageWidth, _ := r.pdf.GetPageSize()
		var margins Margins
		margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
		width := pageWidth - margins.Left - margins.Right

		r.backgroundColor(props.BackgroundColor)

		lineHeight := props.FontSize * 0.5
		if props.LineHeight > 0 {
			lineHeight = props.LineHeight
		}

		r.pdf.MultiCell(width, lineHeight, text, "", align, true)
	} else {
		lineHeight := props.FontSize * 0.5
		if props.LineHeight > 0 {
			lineHeight = props.LineHeight
		}
		r.pdf.MultiCell(0, lineHeight, text, "", align, false)
	}

	r.textColor("")
	r.backgroundColor("")

	if props.MarginBottom > 0 {
		r.pdf.Ln(props.MarginBottom)
	}

	return nil
}

func (r *Renderer) substituteVariables(text string) string {
	return varRe.ReplaceAllStringFunc(text, func(match string) string {
		varName := strings.Trim(match, "{}")

		switch varName {
		case "page":
			return fmt.Sprintf("%d", r.pdf.PageNo())
		case "totalPages":
			return "{nb}"
		}

		val, exists := r.context.GetNested(varName)
		if exists {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}
