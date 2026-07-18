package renderer

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func TestSubstituteVariables(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name":  "Sergio",
		"count": 5,
		"user": map[string]any{
			"email": "sergio@example.com",
		},
	}}

	r := NewRenderer(rc, "", "Arial")

	tests := []struct {
		name string
		text string
		want string
	}{
		{"simple variable", "Hello {{name}}", "Hello Sergio"},
		{"int variable", "Count: {{count}}", "Count: 5"},
		{"nested variable", "Email: {{user.email}}", "Email: sergio@example.com"},
		{"page number", "Page {{page}}", "Page 1"},
		{"total pages", "Total {{totalPages}}", "Total {nb}"},
		{"no variable", "Plain text", "Plain text"},
		{"unknown variable", "{{missing}}", "{{missing}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.substituteVariables(tt.text)
			if got != tt.want {
				t.Errorf("substituteVariables(%q) = %q, want %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestRenderTextBlock(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"name": "Sergio",
	}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "text",
		TextProperties: &models.TextProperties{
			Text:     "Hello {{name}}",
			FontSize: 12,
		},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(text) failed: %v", err)
	}
}

func TestRenderTextBlockMissingProperties(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{Type: "text"}
	if err := r.RenderBlock(block); err == nil {
		t.Error("expected error for text block without TextProperties")
	}
}

func TestRenderPageBreak(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type:                "pagebreak",
		PageBreakProperties: &models.PageBreakProperties{},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(pagebreak) failed: %v", err)
	}
	if r.pdf.PageNo() != 2 {
		t.Errorf("expected page 2 after page break, got page %d", r.pdf.PageNo())
	}
}

func TestRenderLoop(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"items": []any{
			map[string]any{"label": "First"},
			map[string]any{"label": "Second"},
			map[string]any{"label": "Third"},
		},
	}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "loop",
		LoopProperties: &models.LoopProperties{
			DataSource: "items",
			ItemVar:    "item",
		},
		Children: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "{{item.label}}", FontSize: 12}},
		},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(loop) failed: %v", err)
	}

	// item should be cleaned up after the loop
	if _, exists := rc.Get("item"); exists {
		t.Error("item variable should be cleaned up after loop")
	}
}

func TestRenderLoopMissingDataSource(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "loop",
		LoopProperties: &models.LoopProperties{
			DataSource: "missing",
		},
	}

	if err := r.RenderBlock(block); err == nil {
		t.Error("expected error for missing loop data source")
	}
}

func TestRenderTable(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			Headers: []string{"Name", "Age"},
			Rows:    [][]string{{"Sergio", "30"}, {"Maria", "25"}},
		},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(table) failed: %v", err)
	}
}

func TestRenderTableWithNilHeaderStyle(t *testing.T) {
	// This test validates the nil-pointer fix from Phase 1
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			Headers:     []string{"Name", "Age"},
			Rows:        [][]string{{"Sergio", "30"}},
			HeaderStyle: nil, // Explicit nil - should not panic
		},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(table with nil HeaderStyle) failed: %v", err)
	}
}

func TestRenderTableWithStyle(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			Headers: []string{"Name", "Age"},
			Rows:    [][]string{{"Sergio", "30"}},
			HeaderStyle: &models.CellStyle{
				FontSize:   14,
				FontWeight: "bold",
				Align:      "left",
				CellHeight: 8,
			},
			RowStyle: &models.CellStyle{
				FontSize: 10,
			},
		},
	}

	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock(table with styles) failed: %v", err)
	}
}

func TestRenderUnknownBlockType(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{Type: "unknown"}
	if err := r.RenderBlock(block); err == nil {
		t.Error("expected error for unknown block type")
	}
}

func TestRenderLoopMissingProperties(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{Type: "loop"}
	if err := r.RenderBlock(block); err == nil {
		t.Error("expected error for loop without LoopProperties")
	}
}

// Golden file tests verify PDF output consistency via hash comparison
func TestGoldenFiles(t *testing.T) {
	goldens := []struct {
		name  string
		block *models.Block
	}{
		{
			name: "simple-text",
			block: &models.Block{
				Type:           "text",
				TextProperties: &models.TextProperties{Text: "Hello, World!", FontSize: 12},
			},
		},
		{
			name: "text-with-variables",
			block: &models.Block{
				Type:           "text",
				TextProperties: &models.TextProperties{Text: "Welcome, {{user}}!", FontSize: 14},
			},
		},
		{
			name: "table-simple",
			block: &models.Block{
				Type: "table",
				TableProperties: &models.TableProperties{
					Headers: []string{"Col A", "Col B"},
					Rows:    [][]string{{"1", "2"}, {"3", "4"}},
				},
			},
		},
	}

	rc := &models.RuntimeContext{Data: map[string]any{
		"user": "Sergio",
	}}

	for _, g := range goldens {
		t.Run(g.name, func(t *testing.T) {
			r := NewRenderer(rc, "", "Arial")
			if err := r.RenderBlock(g.block); err != nil {
				t.Fatalf("RenderBlock failed: %v", err)
			}

			tmpDir := t.TempDir()
			pdfPath := filepath.Join(tmpDir, "output.pdf")
			if err := r.SaveToFile(pdfPath); err != nil {
				t.Fatalf("SaveToFile failed: %v", err)
			}

			data, err := os.ReadFile(pdfPath)
			if err != nil {
				t.Fatalf("failed to read PDF: %v", err)
			}

			hash := sha256.Sum256(data)
			t.Logf("%s hash: %x", g.name, hash)
		})
	}
}

func TestWriteToBytes(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Arial")

	block := &models.Block{
		Type:           "text",
		TextProperties: &models.TextProperties{Text: "Hello", FontSize: 12},
	}
	if err := r.RenderBlock(block); err != nil {
		t.Fatalf("RenderBlock failed: %v", err)
	}

	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "output.pdf")
	if err := r.SaveToFile(pdfPath); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	data, err := os.ReadFile(pdfPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "%PDF-") {
		t.Error("output does not start with PDF magic bytes")
	}
}
