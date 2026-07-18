package renderer

import (
	"fmt"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

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
