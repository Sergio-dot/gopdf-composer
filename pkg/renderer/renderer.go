package renderer

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
	"github.com/phpdave11/gofpdf"
)

type Renderer struct {
	pdf         *gofpdf.Fpdf
	context     *models.RuntimeContext
	defaultFont string
}

type Margins struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

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

var varRe = regexp.MustCompile(`\{\{([\w.]+)\}\}`)

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

func (r *Renderer) renderImage(block *models.Block) error {
	if block.ImageProperties == nil {
		return fmt.Errorf("image block missing imageProperties")
	}

	props := block.ImageProperties

	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	baseX := r.pdf.GetX()
	baseY := r.pdf.GetY()

	x := baseX
	y := baseY

	if props.Align != "" {
		pageWidth, _ := r.pdf.GetPageSize()
		var margins Margins
		margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
		availableWidth := pageWidth - margins.Left - margins.Right

		switch props.Align {
		case "center":
			x = margins.Left + (availableWidth-props.Width)/2
		case "right":
			x = pageWidth - margins.Right - props.Width
		case "left":
			x = margins.Left
		}
	}

	imgX := x + props.OffsetX
	imgY := y + props.OffsetY

	opts := gofpdf.ImageOptions{ImageType: "", ReadDpi: true}
	r.pdf.ImageOptions(props.Path, imgX, imgY, props.Width, props.Height, false, opts, 0, "")

	renderedHeight := props.Height
	if renderedHeight == 0 {
		info := r.pdf.GetImageInfo(props.Path)
		if info != nil && info.Width() > 0 {
			renderedHeight = props.Width * (info.Height() / info.Width())
		} else {
			renderedHeight = props.Width * 0.75
		}
	}

	finalY := imgY + renderedHeight
	r.pdf.SetY(finalY)

	if props.MarginBottom > 0 {
		r.pdf.Ln(props.MarginBottom)
	}

	return nil
}

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

func (r *Renderer) textColor(hexColor string) {
	if hexColor == "" {
		r.pdf.SetTextColor(0, 0, 0) // default black
		return
	}
	rgb := hexToRGB(hexColor)
	r.pdf.SetTextColor(rgb[0], rgb[1], rgb[2])
}

func (r *Renderer) backgroundColor(hexColor string) {
	if hexColor == "" {
		r.pdf.SetFillColor(255, 255, 255)
		return
	}
	rgb := hexToRGB(hexColor)
	r.pdf.SetFillColor(rgb[0], rgb[1], rgb[2])
}

func (r *Renderer) drawColor(hexColor string) {
	if hexColor == "" {
		r.pdf.SetDrawColor(0, 0, 0)
		return
	}
	rgb := hexToRGB(hexColor)
	r.pdf.SetDrawColor(rgb[0], rgb[1], rgb[2])
}

func hexToRGB(hex string) []int {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return []int{0, 0, 0} // default black
	}

	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return []int{r, g, b}
}

func (r *Renderer) renderContainer(block *models.Block) error {
	if block.Direction == "row" {
		return r.renderRowContainer(block)
	}
	return r.renderColumnContainer(block)
}

// Row container places childrens side by side
func (r *Renderer) renderRowContainer(block *models.Block) error {
	if block.MarginTop != nil && *block.MarginTop > 0 {
		r.pdf.Ln(*block.MarginTop)
	}

	pageWidth, _ := r.pdf.GetPageSize()
	var margins Margins
	margins.Left, margins.Top, margins.Right, _ = r.pdf.GetMargins()
	availableWidth := pageWidth - margins.Left - margins.Right

	startX := margins.Left
	startY := r.pdf.GetY()

	currentX := startX
	maxHeight := 0.0

	for _, child := range block.Children {
		childWidth := availableWidth / float64(len(block.Children))
		if child.WidthPercent != nil {
			childWidth = availableWidth * (*child.WidthPercent / 100)
		}

		r.pdf.SetXY(currentX, startY)

		oldRightMargin := margins.Right
		newRightMargin := pageWidth - (currentX + childWidth)
		r.pdf.SetMargins(currentX, margins.Top, newRightMargin)

		if err := r.RenderBlock(&child); err != nil {
			r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)
			return err
		}

		r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)

		childHeight := r.pdf.GetY() - startY
		if childHeight > maxHeight {
			maxHeight = childHeight
		}

		currentX += childWidth + block.Gap
	}

	if block.BackgroundColor != "" && maxHeight > 0 {
		r.pdf.SetXY(startX, startY)

		r.backgroundColor(block.BackgroundColor)
		if block.Border {
			r.drawColor(block.BackgroundColor)
			r.pdf.Rect(startX, startY, availableWidth, maxHeight, "D")
		} else {
			r.pdf.Rect(startX, startY, availableWidth, maxHeight, "F")
		}
		r.backgroundColor("")

		currentX = startX
		for _, child := range block.Children {
			childWidth := availableWidth / float64(len(block.Children))
			if child.WidthPercent != nil {
				childWidth = availableWidth * (*child.WidthPercent / 100)
			}

			r.pdf.SetXY(currentX, startY)

			oldRightMargin := margins.Right
			newRightMargin := pageWidth - (currentX + childWidth)
			r.pdf.SetMargins(currentX, margins.Top, newRightMargin)

			if err := r.RenderBlock(&child); err != nil {
				r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)
				return err
			}

			r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)

			currentX += childWidth + block.Gap
		}
	}

	r.pdf.SetXY(margins.Left, startY+maxHeight)

	return nil
}

// Column container stacks children vertically
func (r *Renderer) renderColumnContainer(block *models.Block) error {
	if block.MarginTop != nil && *block.MarginTop > 0 {
		r.pdf.Ln(*block.MarginTop)
	}

	startY := r.pdf.GetY()
	startX := r.pdf.GetX()

	for i, child := range block.Children {
		if err := r.RenderBlock(&child); err != nil {
			return err
		}

		if i < len(block.Children)-1 && block.Gap > 0 {
			r.pdf.Ln(block.Gap)
		}
	}

	endY := r.pdf.GetY()
	totalHeight := endY - startY

	if block.BackgroundColor != "" && totalHeight > 0 {
		pageWidth, _ := r.pdf.GetPageSize()
		var margins Margins
		margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
		availableWidth := pageWidth - margins.Left - margins.Right

		r.pdf.SetXY(startX, startY)

		r.backgroundColor(block.BackgroundColor)
		if block.Border {
			r.pdf.Rect(startX, startY, availableWidth, totalHeight, "D")
		} else {
			r.pdf.Rect(startX, startY, availableWidth, totalHeight, "F")
		}
		r.backgroundColor("")

		r.pdf.SetXY(startX, startY)
		for i, child := range block.Children {
			if err := r.RenderBlock(&child); err != nil {
				return err
			}

			if i < len(block.Children)-1 && block.Gap > 0 {
				r.pdf.Ln(block.Gap)
			}
		}
	}

	return nil
}

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

func (r *Renderer) RenderBlocks(blocks []models.Block) error {
	for _, block := range blocks {
		if err := r.RenderBlock(&block); err != nil {
			return err
		}
	}
	return nil
}

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

func (r *Renderer) RenderHeader(blocks []models.Block) error {
	return r.RenderBlocks(blocks)
}

func (r *Renderer) RenderFooter(blocks []models.Block, offsetFromBottom float64) error {
	r.pdf.SetY(-offsetFromBottom)
	return r.RenderBlocks(blocks)
}

func (r *Renderer) GetPDF() *gofpdf.Fpdf {
	return r.pdf
}

func (r *Renderer) GetContext() *models.RuntimeContext {
	return r.context
}

func (r *Renderer) renderLoop(block *models.Block) error {
	if block.LoopProperties == nil {
		return fmt.Errorf("loop block missing loopProperties")
	}

	props := block.LoopProperties

	dataSource, exists := r.context.Get(props.DataSource)
	if !exists {
		return fmt.Errorf("loop dataSource not found in context: %s", props.DataSource)
	}

	items, ok := dataSource.([]any)
	if !ok {
		return fmt.Errorf("loop dataSource is not an array: %s", props.DataSource)
	}

	itemVar := props.ItemVar
	if itemVar == "" {
		itemVar = "item"
	}

	for _, item := range items {
		r.context.Set(itemVar, item)
		for _, child := range block.Children {
			if err := r.RenderBlock(&child); err != nil {
				r.context.Delete(itemVar)
				return err
			}
		}
		r.context.Delete(itemVar)
	}

	return nil
}

func (r *Renderer) SaveToFile(path string) error {
	return r.pdf.OutputFileAndClose(path)
}

func (r *Renderer) WriteTo(w io.Writer) (int64, error) {
	return 0, r.pdf.Output(w)
}
