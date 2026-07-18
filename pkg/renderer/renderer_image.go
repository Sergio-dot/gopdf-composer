package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
	"github.com/phpdave11/gofpdf"
)

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
