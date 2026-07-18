package renderer

import "github.com/Sergio-dot/gopdf-composer/pkg/models"

func (r *Renderer) renderContainer(block *models.Block) error {
	if block.Direction == "row" {
		return r.renderRowContainer(block)
	}
	return r.renderColumnContainer(block)
}

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
