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

func NewRenderer(runtimeCtx *models.RuntimeContext, fontDir, defaultFont string) *Renderer {
	pdf := gofpdf.New("P", "mm", "A4", fontDir)

	if defaultFont == "" {
		defaultFont = "Arial"
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
	default:
		return fmt.Errorf("unknown block type: %s", block.Type)
	}
}

func (r *Renderer) renderText(block *models.Block) error {
	if block.TextProperties == nil {
		return fmt.Errorf("text block missing textProperties")
	}

	props := block.TextProperties

	// top margin
	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	// font
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

	// text color
	r.textColor(props.FontColor)

	// substitute variables
	text := r.substituteVariables(props.Text) // TODO: maybe is better to use text/template

	// alignment
	align := "L"
	switch props.Align {
	case "center":
		align = "C"
	case "right":
		align = "R"
	}

	if props.BackgroundColor != "" { // with bg color
		pageWidth, _ := r.pdf.GetPageSize()
		var margins Margins
		margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
		width := pageWidth - margins.Left - margins.Right

		// background color
		r.backgroundColor(props.BackgroundColor)

		// line height
		lineHeight := props.FontSize * 0.5
		if props.LineHeight > 0 {
			lineHeight = props.LineHeight
		}

		r.pdf.MultiCell(width, lineHeight, text, "", align, true)
	} else { // no bg color
		lineHeight := props.FontSize * 0.5
		if props.LineHeight > 0 {
			lineHeight = props.LineHeight
		}
		r.pdf.MultiCell(0, lineHeight, text, "", align, false)
	}

	// reset colors
	r.textColor("")
	r.backgroundColor("")

	// bottom margin
	if props.MarginBottom > 0 {
		r.pdf.Ln(props.MarginBottom)
	}

	return nil
}

func (r *Renderer) substituteVariables(text string) string {
	re := regexp.MustCompile(`\{\{([\w.]+)\}\}`)

	return re.ReplaceAllStringFunc(text, func(match string) string {
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

	// Handle top margin
	if props.MarginTop > 0 {
		r.pdf.Ln(props.MarginTop)
	}

	// Get current position
	baseX := r.pdf.GetX()
	baseY := r.pdf.GetY()

	x := baseX
	y := baseY

	// Handle alignment
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

	// Apply offsets
	imgX := x + props.OffsetX
	imgY := y + props.OffsetY

	// Render image
	opts := gofpdf.ImageOptions{ImageType: "", ReadDpi: true}
	r.pdf.ImageOptions(props.Path, imgX, imgY, props.Width, props.Height, false, opts, 0, "")

	// Move cursor below image
	renderedHeight := props.Height
	if renderedHeight == 0 {
		// Attempt to get image info to find aspect ratio if height is auto
		info := r.pdf.GetImageInfo(props.Path)
		if info != nil && info.Width() > 0 {
			renderedHeight = props.Width * (info.Height() / info.Width())
		} else {
			// Fallback if info not available
			renderedHeight = props.Width * 0.75
		}
	}

	// Final Y is the bottom of the image
	finalY := imgY + renderedHeight
	r.pdf.SetY(finalY)

	// Handle bottom margin
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

	// Calculate column widths
	pageWidth, _ := r.pdf.GetPageSize()
	var margins Margins
	margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
	availableWidth := pageWidth - margins.Left - margins.Right
	colWidth := availableWidth / float64(len(props.Headers))

	// Render headers
	if props.HeaderStyle != nil {
		r.applyTableCellStyle(props.HeaderStyle)
	} else {
		r.textColor("")
		r.backgroundColor("")
		r.pdf.SetFont(r.defaultFont, "B", 10)
	}

	align := "C"
	switch props.HeaderStyle.Align {
	case "left":
		align = "L"
	case "right":
		align = "R"
	}

	cellHeight := 5.0
	if props.HeaderStyle.CellHeight > 0 {
		cellHeight = props.HeaderStyle.CellHeight
	}

	for _, header := range props.Headers {
		r.pdf.CellFormat(colWidth, cellHeight, header, "", 0, align, true, 0, "")
	}
	r.pdf.Ln(-1)

	// Render rows
	if props.RowStyle != nil {
		r.applyTableCellStyle(props.RowStyle)
	} else {
		r.textColor("")
		r.backgroundColor("")
		r.pdf.SetFont(r.defaultFont, "", 9)
	}

	if props.RowStyle.CellHeight > 0 {
		cellHeight = props.RowStyle.CellHeight
	}

	for _, row := range props.Rows {
		for _, cell := range row {
			cell = r.substituteVariables(cell)
			r.pdf.CellFormat(colWidth, cellHeight, cell, "", 0, "L", false, 0, "")
		}
		r.pdf.Ln(-1)
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

		// apply headers colors
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

	// 1: calculate max height by rendering children
	currentX := startX
	maxHeight := 0.0

	for _, child := range block.Children {
		// Calculate child width
		childWidth := availableWidth / float64(len(block.Children))
		if child.WidthPercent != nil {
			childWidth = availableWidth * (*child.WidthPercent / 100)
		}

		// Set position for current child
		r.pdf.SetXY(currentX, startY)

		// Save the current margins and set temporary right margin for this column
		oldRightMargin := margins.Right
		newRightMargin := pageWidth - (currentX + childWidth)
		r.pdf.SetMargins(currentX, margins.Top, newRightMargin)

		// Render child
		if err := r.RenderBlock(&child); err != nil {
			// Restore margins before returning error
			r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)
			return err
		}

		// Restore original margins
		r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)

		// Update max height
		childHeight := r.pdf.GetY() - startY
		if childHeight > maxHeight {
			maxHeight = childHeight
		}

		// Move to next col
		currentX += childWidth + block.Gap
	}

	// draw background BEHIND the content if specified
	if block.BackgroundColor != "" && maxHeight > 0 {
		// Go back to start position to draw background
		r.pdf.SetXY(startX, startY)

		// draw rectangle
		r.backgroundColor(block.BackgroundColor)
		if block.Border {
			r.drawColor(block.BackgroundColor)
			r.pdf.Rect(startX, startY, availableWidth, maxHeight, "D")
		} else {
			r.pdf.Rect(startX, startY, availableWidth, maxHeight, "F")
		}
		r.backgroundColor("")

		// 2: re-render children on top of background
		currentX = startX
		for _, child := range block.Children {
			// Calculate child width
			childWidth := availableWidth / float64(len(block.Children))
			if child.WidthPercent != nil {
				childWidth = availableWidth * (*child.WidthPercent / 100)
			}

			// Set position for current child
			r.pdf.SetXY(currentX, startY)

			// Save the current margins and set temporary right margin for this column
			oldRightMargin := margins.Right
			newRightMargin := pageWidth - (currentX + childWidth)
			r.pdf.SetMargins(currentX, margins.Top, newRightMargin)

			// Render child
			if err := r.RenderBlock(&child); err != nil {
				// Restore margins before returning error
				r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)
				return err
			}

			// Restore original margins
			r.pdf.SetMargins(margins.Left, margins.Top, oldRightMargin)

			// Move to next col
			currentX += childWidth + block.Gap
		}
	}

	// Move cursor below the entire row
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

	// 1: render children to calculate total height
	for i, child := range block.Children {
		if err := r.RenderBlock(&child); err != nil {
			return err
		}

		// add gap
		if i < len(block.Children)-1 && block.Gap > 0 {
			r.pdf.Ln(block.Gap)
		}
	}

	endY := r.pdf.GetY()
	totalHeight := endY - startY

	// draw background and re-render if background color specified
	if block.BackgroundColor != "" && totalHeight > 0 {
		pageWidth, _ := r.pdf.GetPageSize()
		var margins Margins
		margins.Left, _, margins.Right, _ = r.pdf.GetMargins()
		availableWidth := pageWidth - margins.Left - margins.Right

		// Go back to start to draw background
		r.pdf.SetXY(startX, startY)

		// draw rectangle
		r.backgroundColor(block.BackgroundColor)
		if block.Border {
			r.pdf.Rect(startX, startY, availableWidth, totalHeight, "D")
		} else {
			r.pdf.Rect(startX, startY, availableWidth, totalHeight, "F")
		}
		r.backgroundColor("")

		// 2: re-render children on top of background
		r.pdf.SetXY(startX, startY)
		for i, child := range block.Children {
			if err := r.RenderBlock(&child); err != nil {
				return err
			}

			// add gap
			if i < len(block.Children)-1 && block.Gap > 0 {
				r.pdf.Ln(block.Gap)
			}
		}
	}

	return nil
}

func (r *Renderer) renderPageBreak(block *models.Block) error {
	r.pdf.AddPage()
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
