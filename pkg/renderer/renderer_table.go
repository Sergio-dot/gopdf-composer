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

	headers := props.Headers
	if props.HeadersDataSource != "" {
		data, ok := r.resolveContext(props.HeadersDataSource)
		if !ok {
			return fmt.Errorf("table headersDataSource not found in context: %s", props.HeadersDataSource)
		}
		hs, ok := data.([]any)
		if !ok {
			return fmt.Errorf("table headersDataSource is not an array: %s", props.HeadersDataSource)
		}
		headers = make([]string, len(hs))
		for i, h := range hs {
			headers[i] = fmt.Sprintf("%v", h)
		}
	}

	pageWidth, _ := r.pdf.GetPageSize()
	var margins Margins
	margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
	availableWidth := pageWidth - margins.Left - margins.Right

	colWidths := make([]float64, len(headers))
	if len(props.ColumnWidths) == len(headers) {
		var total float64
		for _, w := range props.ColumnWidths {
			total += w
		}
		if total > 0 {
			for i, w := range props.ColumnWidths {
				colWidths[i] = (w / total) * availableWidth
			}
		} else {
			for i := range colWidths {
				colWidths[i] = availableWidth / float64(len(headers))
			}
		}
	} else {
		for i := range colWidths {
			colWidths[i] = availableWidth / float64(len(headers))
		}
	}

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

	headerStartY := r.pdf.GetY()
	maxHeaderLines := 1
	for i, header := range headers {
		lines := len(r.pdf.SplitLines([]byte(header), colWidths[i]))
		if lines > maxHeaderLines {
			maxHeaderLines = lines
		}
	}
	headerRowHeight := float64(maxHeaderLines) * cellHeight

	rectStyle := "F"
	if props.Border {
		rectStyle = "FD"
	}

	x := margins.Left
	for i := range headers {
		r.pdf.Rect(x, headerStartY, colWidths[i], headerRowHeight, rectStyle)
		x += colWidths[i]
	}

	x = margins.Left
	for i, header := range headers {
		r.pdf.SetXY(x, headerStartY)
		r.pdf.MultiCell(colWidths[i], cellHeight, header, "", align, false)
		x += colWidths[i]
	}
	r.pdf.SetY(headerStartY + headerRowHeight)

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

	rowAlign := "L"
	if props.RowStyle != nil {
		switch props.RowStyle.Align {
		case "left":
			rowAlign = "L"
		case "right":
			rowAlign = "R"
		case "center":
			rowAlign = "C"
		}
	}

	if props.RowsDataSource != "" {
		data, ok := r.resolveContext(props.RowsDataSource)
		if !ok {
			return fmt.Errorf("table rowsDataSource not found in context: %s", props.RowsDataSource)
		}
		items, ok := data.([]any)
		if !ok {
			return fmt.Errorf("table rowsDataSource is not an array: %s", props.RowsDataSource)
		}
		if props.HeadersDataSource != "" {
			cellsField := props.CellsField
			if cellsField == "" {
				cellsField = "cells"
			}
			for _, item := range items {
				cells := rowCellsFromItem(item, cellsField)
				r.renderRowCells(colWidths, cellHeight, cells, rowAlign, props.Border, props.StrikeEmpty)
			}
		} else {
			if len(props.Rows) == 0 {
				return fmt.Errorf("table with rowsDataSource must have at least one template row")
			}
			templateRow := props.Rows[0]
			for _, item := range items {
				r.context.Set("item", item)
				cells := make([]string, len(templateRow))
				for i, cell := range templateRow {
					cells[i] = r.substituteVariables(cell)
				}
				r.renderRowCells(colWidths, cellHeight, cells, rowAlign, props.Border, props.StrikeEmpty)
				r.context.Delete("item")
			}
		}
	} else {
		for _, row := range props.Rows {
			cells := make([]string, len(row))
			for i, cell := range row {
				cells[i] = r.substituteVariables(cell)
			}
			r.renderRowCells(colWidths, cellHeight, cells, rowAlign, props.Border, props.StrikeEmpty)
		}
	}

	return nil
}

// resolveContext looks up a value by top-level key first, then by nested
// dot path (so loop items such as "item.rows" are supported).
func (r *Renderer) resolveContext(key string) (any, bool) {
	if v, ok := r.context.Get(key); ok {
		return v, true
	}
	if v, ok := r.context.GetNested(key); ok {
		return v, true
	}
	return nil, false
}

// rowCellsFromItem extracts the "cells" array from a dynamic row item.
func rowCellsFromItem(item any, cellsField string) []string {
	m, ok := item.(map[string]any)
	if !ok {
		return nil
	}
	raw, ok := m[cellsField]
	if !ok {
		return nil
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	cells := make([]string, len(arr))
	for i, v := range arr {
		cells[i] = fmt.Sprintf("%v", v)
	}
	return cells
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

func (r *Renderer) renderRowCells(colWidths []float64, lineHt float64, cells []string, align string, border bool, strikeEmpty bool) {
	marginLeft, _, _, _ := r.pdf.GetMargins()

	maxLines := 0
	for i, text := range cells {
		var lines [][]byte
		func() {
			defer func() {
				if recover() != nil {
					lines = [][]byte{[]byte(text)}
				}
			}()
			lines = r.pdf.SplitLines([]byte(text), colWidths[i])
		}()
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}
	if maxLines == 0 {
		maxLines = 1
	}

	neededHt := float64(maxLines) * lineHt
	if !r.fitsOnPage(neededHt) {
		r.pdf.AddPage()
	}

	var colX float64 = marginLeft
	colStartX := make([]float64, len(colWidths))
	for i, w := range colWidths {
		colStartX[i] = colX
		colX += w
	}

	startY := r.pdf.GetY()
	maxY := startY

	borderStr := ""
	if border {
		borderStr = "1"
	}

	for i, text := range cells {
		r.pdf.SetY(startY)
		r.pdf.SetX(colStartX[i])
		r.pdf.MultiCell(colWidths[i], lineHt, text, borderStr, align, false)
		if strikeEmpty && text == "" {
			r.pdf.Line(colStartX[i], startY, colStartX[i]+colWidths[i], startY+lineHt)
		}
		if cy := r.pdf.GetY(); cy > maxY {
			maxY = cy
		}
	}

	r.pdf.SetY(maxY)
}

func (r *Renderer) fitsOnPage(neededHt float64) bool {
	_, pageHt := r.pdf.GetPageSize()
	_, _, _, bottomMargin := r.pdf.GetMargins()
	auto, autoMargin := r.pdf.GetAutoPageBreak()
	if auto {
		bottomMargin = autoMargin
	}
	return r.pdf.GetY()+neededHt <= pageHt-bottomMargin
}
