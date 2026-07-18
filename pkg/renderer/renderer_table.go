package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func (r *Renderer) renderTable(block *models.Block) error {
	if block.TableProperties == nil {
		return fmt.Errorf("table block missing tableProperties")
	}

	props := block.TableProperties

	pageWidth, _ := r.pdf.GetPageSize()
	var margins Margins
	margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
	availableWidth := pageWidth - margins.Left - margins.Right
	colWidth := availableWidth / float64(len(props.Headers))

	align := "C"
	cellHeight := 5.0

	if props.HeaderStyle != nil {
		r.applyTableCellStyle(props.HeaderStyle)
		switch props.HeaderStyle.Align {
		case "left":
			align = "L"
		case "right":
			align = "R"
		}
		if props.HeaderStyle.CellHeight > 0 {
			cellHeight = props.HeaderStyle.CellHeight
		}
	} else {
		r.textColor("")
		r.backgroundColor("")
		r.pdf.SetFont(r.defaultFont, "B", 10)
	}

	for _, header := range props.Headers {
		r.pdf.CellFormat(colWidth, cellHeight, header, "", 0, align, true, 0, "")
	}
	r.pdf.Ln(-1)

	if props.RowStyle != nil {
		r.applyTableCellStyle(props.RowStyle)
	} else {
		r.textColor("")
		r.backgroundColor("")
		r.pdf.SetFont(r.defaultFont, "", 9)
	}

	if props.RowStyle != nil && props.RowStyle.CellHeight > 0 {
		cellHeight = props.RowStyle.CellHeight
	}

	if props.RowsDataSource != "" {
		data, exists := r.context.Get(props.RowsDataSource)
		if !exists {
			return fmt.Errorf("table rowsDataSource not found in context: %s", props.RowsDataSource)
		}
		items, ok := data.([]any)
		if !ok {
			return fmt.Errorf("table rowsDataSource is not an array: %s", props.RowsDataSource)
		}
		if len(props.Rows) == 0 {
			return fmt.Errorf("table with rowsDataSource must have at least one template row")
		}
		templateRow := props.Rows[0]
		for _, item := range items {
			r.context.Set("item", item)
			for _, cell := range templateRow {
				cell = r.substituteVariables(cell)
				r.pdf.CellFormat(colWidth, cellHeight, cell, "", 0, "L", false, 0, "")
			}
			r.pdf.Ln(-1)
			r.context.Delete("item")
		}
	} else {
		for _, row := range props.Rows {
			for _, cell := range row {
				cell = r.substituteVariables(cell)
				r.pdf.CellFormat(colWidth, cellHeight, cell, "", 0, "L", false, 0, "")
			}
			r.pdf.Ln(-1)
		}
	}

	return nil
}

func (r *Renderer) applyTableCellStyle(style *models.CellStyle) {
	if style.FontSize > 0 {
		fontFamily := r.defaultFont

		fontStyle := ""
		switch style.FontWeight {
		case "bold":
			fontStyle = "B"
		case "italic":
			fontStyle = "I"
		}

		r.textColor(style.FontColor)
		r.backgroundColor(style.BackgroundColor)

		r.pdf.SetFont(fontFamily, fontStyle, style.FontSize)
	}
}
