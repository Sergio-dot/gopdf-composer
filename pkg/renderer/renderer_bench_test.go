package renderer

import (
	"testing"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func BenchmarkSubstituteVariables(b *testing.B) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name": "Sergio",
		"user": map[string]any{"email": "sergio@example.com"},
	}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	text := "Hello {{name}}, your email is {{user.email}}. Welcome to page {{page}} of {{totalPages}}."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.substituteVariables(text)
	}
}

func BenchmarkRenderTextBlock(b *testing.B) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name": "Sergio",
	}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	block := &models.Block{
		Type: "text",
		TextProperties: &models.TextProperties{
			Text:     "Hello {{name}}, this is a benchmark text block.",
			FontSize: 12,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderText(block)
	}
}

func BenchmarkRenderTableBlock(b *testing.B) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	rows := make([][]string, 20)
	for i := 0; i < 20; i++ {
		rows[i] = []string{"Name", "Value", "Status"}
	}

	block := &models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			Headers: []string{"Col A", "Col B", "Col C"},
			Rows:    rows,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderTable(block)
	}
}

func BenchmarkRenderRowContainer(b *testing.B) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	block := &models.Block{
		Type:      "container",
		Direction: "row",
		Children: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Column 1", FontSize: 10}},
			{Type: "text", TextProperties: &models.TextProperties{Text: "Column 2", FontSize: 10}},
			{Type: "text", TextProperties: &models.TextProperties{Text: "Column 3", FontSize: 10}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderRowContainer(block)
	}
}

func BenchmarkRenderRowContainerWithBackground(b *testing.B) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	block := &models.Block{
		Type:            "container",
		Direction:       "row",
		BackgroundColor: "#f0f0f0",
		Children: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "Column 1", FontSize: 10}},
			{Type: "text", TextProperties: &models.TextProperties{Text: "Column 2", FontSize: 10}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderRowContainer(block)
	}
}

func BenchmarkRenderLoop(b *testing.B) {
	items := make([]any, 10)
	for i := 0; i < 10; i++ {
		items[i] = map[string]any{"label": "Item"}
	}
	rc := &models.RuntimeContext{Data: map[string]any{"items": items}}
	r := NewRenderer(rc, "", "Arial", "", "", nil)

	block := &models.Block{
		Type: "loop",
		LoopProperties: &models.LoopProperties{
			DataSource: "items",
			ItemVar:    "item",
		},
		Children: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "{{item.label}}", FontSize: 10}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderLoop(block)
	}
}
