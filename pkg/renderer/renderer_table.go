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

	colWidths := make([]float64, len(props.Headers))
	if len(props.ColumnWidths) == len(props.Headers) {
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
				colWidths[i] = availableWidth / float64(len(props.Headers))
			}
		}
	} else {
		for i := range colWidths {
			colWidths[i] = availableWidth / float64(len(props.Headers))
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

	for i, header := range props.Headers {
		r.pdf.CellFormat(colWidths[i], cellHeight, header, "", 0, align, true, 0, "")
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
			cells := make([]string, len(templateRow))
			for i, cell := range templateRow {
				cells[i] = r.substituteVariables(cell)
			}
			r.renderRowCells(colWidths, cellHeight, cells, "L")
			r.context.Delete("item")
		}
	} else {
		for _, row := range props.Rows {
			cells := make([]string, len(row))
			for i, cell := range row {
				cells[i] = r.substituteVariables(cell)
			}
			r.renderRowCells(colWidths, cellHeight, cells, "L")
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

func (r *Renderer) renderRowCells(colWidths []float64, lineHt float64, cells []string, align string) {
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

	for i, text := range cells {
		r.pdf.SetY(startY)
		r.pdf.SetX(colStartX[i])
		r.pdf.MultiCell(colWidths[i], lineHt, text, "", align, false)
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
