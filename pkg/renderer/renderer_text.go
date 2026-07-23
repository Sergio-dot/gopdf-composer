package renderer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

var varRe = regexp.MustCompile(`\{\{([\w.]+)\}\}`)

func resolveFontStyle(weight string) string {
	switch weight {
	case "bold":
		return "B"
	case "italic":
		return "I"
	case "bold-italic":
		return "BI"
	default:
		return ""
	}
}

func (r *Renderer) renderText(block *models.Block) error {
	if block.TextProperties == nil {
		return fmt.Errorf("text block missing textProperties")
	}

	props := block.TextProperties

	if len(props.Spans) > 0 {
		return r.renderTextSpans(props)
	}

	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	fontFamily := props.FontFamily
	if fontFamily == "" {
		fontFamily = r.defaultFont
	}

	fontWeight := resolveFontStyle(props.FontWeight)
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

func (r *Renderer) renderTextSpans(props *models.TextProperties) error {
	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	baseFontFamily := props.FontFamily
	if baseFontFamily == "" {
		baseFontFamily = r.defaultFont
	}
	baseFontSize := props.FontSize

	lineHeight := props.LineHeight
	if lineHeight <= 0 {
		maxFontSize := baseFontSize
		for _, span := range props.Spans {
			if span.FontSize > maxFontSize {
				maxFontSize = span.FontSize
			}
		}
		lineHeight = maxFontSize * 0.5
	}

	totalWidth := 0.0
	for _, span := range props.Spans {
		ff := span.FontFamily
		if ff == "" {
			ff = baseFontFamily
		}
		fs := span.FontSize
		if fs == 0 {
			fs = baseFontSize
		}
		style := resolveFontStyle(span.FontWeight)
		r.pdf.SetFont(ff, style, fs)
		totalWidth += r.pdf.GetStringWidth(r.substituteVariables(span.Text))
	}

	align := "L"
	switch props.Align {
	case "center":
		align = "C"
	case "right":
		align = "R"
	}
	if align != "L" {
		pageWidth, _ := r.pdf.GetPageSize()
		leftMargin, _, rightMargin, _ := r.pdf.GetMargins()
		availableWidth := pageWidth - leftMargin - rightMargin
		offset := availableWidth - totalWidth
		if align == "C" {
			offset /= 2
		}
		if offset > 0 {
			r.pdf.SetX(leftMargin + offset)
		}
	}

	for _, span := range props.Spans {
		ff := span.FontFamily
		if ff == "" {
			ff = baseFontFamily
		}
		fs := span.FontSize
		if fs == 0 {
			fs = baseFontSize
		}
		style := resolveFontStyle(span.FontWeight)
		r.pdf.SetFont(ff, style, fs)
		r.textColor(span.FontColor)
		r.pdf.Write(lineHeight, r.substituteVariables(span.Text))
	}

	r.textColor("")

	leftMargin, _, _, _ := r.pdf.GetMargins()
	if r.pdf.GetX() != leftMargin {
		r.pdf.Ln(lineHeight)
	}

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
