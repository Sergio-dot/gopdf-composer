package renderer

import (
	"fmt"
	"strings"
)

func (r *Renderer) textColor(hexColor string) {
	if hexColor == "" {
		r.pdf.SetTextColor(0, 0, 0)
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
		return []int{0, 0, 0}
	}

	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return []int{r, g, b}
}
