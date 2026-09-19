package renderer

import (
	"bytes"
	"testing"

	"github.com/Sergio-dot/gopdf-composer/pkg/models"
)

func TestRenderDynamicTableWithBordersAndStrike(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"headers": []any{"Giorno", "Frigo 1", "Freezer", "FIRMA"},
		"rows": []any{
			map[string]any{"cells": []any{1, "4.2", "-18.0", ""}},
			map[string]any{"cells": []any{2, "", "-18.5", ""}},
		},
	}}

	r := NewRenderer(rc, "", "Helvetica", "", "", nil)

	table := models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			HeadersDataSource: "headers",
			RowsDataSource:    "rows",
			Border:            true,
			StrikeEmpty:       true,
			HeaderStyle:       &models.CellStyle{CellHeight: 6, FontSize: 8},
			RowStyle:          &models.CellStyle{CellHeight: 6, FontSize: 8},
		},
	}

	if err := r.RenderBlock(&table); err != nil {
		t.Fatalf("render dynamic table: %v", err)
	}

	var buf bytes.Buffer
	if err := r.pdf.Output(&buf); err != nil {
		t.Fatalf("pdf output: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty pdf output")
	}
}

func TestRenderDynamicTableInsideLoop(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"sheets": []any{
			map[string]any{
				"monthLabel": "08/2026",
				"headers":    []any{"Giorno", "Frigo", "FIRMA"},
				"rows": []any{
					map[string]any{"cells": []any{1, "4.0", ""}},
					map[string]any{"cells": []any{2, "4.1", ""}},
				},
			},
		},
	}}

	r := NewRenderer(rc, "", "Helvetica", "", "", nil)

	loop := models.Block{
		Type:           "loop",
		LoopProperties: &models.LoopProperties{DataSource: "sheets", ItemVar: "item"},
		Children: []models.Block{
			{Type: "text", TextProperties: &models.TextProperties{Text: "3. SCHEDA CONTROLLO TEMPERATURE Mese {{item.monthLabel}}", FontSize: 12}},
			{
				Type: "table",
				TableProperties: &models.TableProperties{
					HeadersDataSource: "item.headers",
					RowsDataSource:    "item.rows",
					Border:            true,
					HeaderStyle:       &models.CellStyle{CellHeight: 6, FontSize: 8},
					RowStyle:          &models.CellStyle{CellHeight: 6, FontSize: 8},
				},
			},
		},
	}

	if err := r.RenderBlock(&loop); err != nil {
		t.Fatalf("render loop table: %v", err)
	}

	var buf bytes.Buffer
	if err := r.pdf.Output(&buf); err != nil {
		t.Fatalf("pdf output: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty pdf output")
	}
}

func TestRenderSignatureBlock(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Helvetica", "", "", nil)

	sig := models.Block{
		Type: "signature",
		SignatureProperties: &models.SignatureProperties{
			Label:          "Il Responsabile dell'autocontrollo",
			SignatureWidth: 120,
			LineWidth:      0.5,
			MarginTop:      10,
		},
	}

	if err := r.RenderBlock(&sig); err != nil {
		t.Fatalf("render signature (inline): %v", err)
	}

	var buf bytes.Buffer
	if err := r.pdf.Output(&buf); err != nil {
		t.Fatalf("pdf output: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty pdf output")
	}
}

func TestRenderMultilineHeaders(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{
		"headers": []any{"DEPOSITO\n(settimanale)", "PIANI LAVORO\n(giornaliera)", "FIRMA"},
		"rows": []any{
			map[string]any{"cells": []any{"", "", ""}},
		},
	}}

	r := NewRenderer(rc, "", "Helvetica", "", "", nil)

	table := models.Block{
		Type: "table",
		TableProperties: &models.TableProperties{
			HeadersDataSource: "headers",
			RowsDataSource:    "rows",
			Border:            true,
			HeaderStyle:       &models.CellStyle{CellHeight: 5, FontSize: 6, Align: "center"},
			RowStyle:          &models.CellStyle{CellHeight: 6, FontSize: 6},
		},
	}

	if err := r.RenderBlock(&table); err != nil {
		t.Fatalf("render multiline-header table: %v", err)
	}

	var buf bytes.Buffer
	if err := r.pdf.Output(&buf); err != nil {
		t.Fatalf("pdf output: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty pdf output")
	}
}

func TestRenderSignatureStackedLayout(t *testing.T) {
	rc := &models.RuntimeContext{Data: map[string]any{}}
	r := NewRenderer(rc, "", "Helvetica", "", "", nil)

	sig := models.Block{
		Type: "signature",
		SignatureProperties: &models.SignatureProperties{
			Label:        "Firma del responsabile",
			Layout:       "stacked",
			FontSize:     11,
			FontWeight:   "bold",
			LineWidth:    0.4,
			LabelLineGap: 12,
		},
	}

	if err := r.RenderBlock(&sig); err != nil {
		t.Fatalf("render signature (below): %v", err)
	}

	var buf bytes.Buffer
	if err := r.pdf.Output(&buf); err != nil {
		t.Fatalf("pdf output: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty pdf output")
	}
}
