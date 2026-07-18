package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRuntimeContextGet(t *testing.T) {
	rc := &RuntimeContext{Data: map[string]any{
		"name": "Sergio",
		"age":  30,
	}}

	if val, ok := rc.Get("name"); !ok || val != "Sergio" {
		t.Errorf("expected 'Sergio', got %v (ok=%v)", val, ok)
	}
	if _, ok := rc.Get("missing"); ok {
		t.Error("expected missing key to not exist")
	}
}

func TestRuntimeContextGetNested(t *testing.T) {
	rc := &RuntimeContext{Data: map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"email": "sergio@example.com",
			},
		},
		"items": []any{"a", "b", "c"},
	}}

	tests := []struct {
		path   string
		want   any
		wantOK bool
	}{
		{"user.profile.email", "sergio@example.com", true},
		{"user.profile.missing", nil, false},
		{"missing.anything", nil, false},
		{"items", []any{"a", "b", "c"}, true},
		{"", nil, false},
	}

	for _, tt := range tests {
		val, ok := rc.GetNested(tt.path)
		if ok != tt.wantOK {
			t.Errorf("GetNested(%q) ok=%v, want %v", tt.path, ok, tt.wantOK)
			continue
		}
		if !reflect.DeepEqual(val, tt.want) {
			t.Errorf("GetNested(%q) = %v, want %v", tt.path, val, tt.want)
		}
	}
}

func TestRuntimeContextSetDelete(t *testing.T) {
	rc := &RuntimeContext{Data: map[string]any{}}
	rc.Set("key", "value")
	if val, ok := rc.Get("key"); !ok || val != "value" {
		t.Errorf("Set/Get failed: got %v", val)
	}
	rc.Delete("key")
	if _, ok := rc.Get("key"); ok {
		t.Error("Delete failed: key still exists")
	}
}

func TestControlFlowJSONRoundtrip(t *testing.T) {
	cf := ControlFlow{
		Document: Document{
			Structure: []Section{
				{Assets: []AssetReference{
					{AssetID: "header", Version: "1"},
					{AssetID: "body", Version: "2", Conditions: &Condition{
						Field: "showBody", Op: "==", Value: true,
					}},
				}},
			},
			HeaderAssets: []AssetReference{
				{AssetID: "page-header", Version: "1"},
			},
			FooterAssets: []AssetReference{
				{AssetID: "page-footer", Version: "1"},
			},
		},
	}

	data, err := json.Marshal(cf)
	if err != nil {
		t.Fatal(err)
	}

	var decoded ControlFlow
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if len(decoded.Document.Structure) != 1 {
		t.Errorf("expected 1 section, got %d", len(decoded.Document.Structure))
	}
	if len(decoded.Document.HeaderAssets) != 1 {
		t.Errorf("expected 1 header asset, got %d", len(decoded.Document.HeaderAssets))
	}
	if decoded.Document.Structure[0].Assets[1].Conditions.Value != true {
		t.Error("condition value was not preserved")
	}
}

func TestBlockJSONRoundtrip(t *testing.T) {
	wp := 80.0
	asset := Asset{
		Blocks: []Block{
			{Type: "text", TextProperties: &TextProperties{Text: "Hello", FontSize: 12}},
			{Type: "pagebreak", PageBreakProperties: &PageBreakProperties{}},
			{Type: "container", Direction: "row", Children: []Block{
				{Type: "text", TextProperties: &TextProperties{Text: "Child", FontSize: 10}, WidthPercent: &wp},
			}},
		},
	}

	data, err := json.Marshal(asset)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Asset
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if len(decoded.Blocks) != 3 {
		t.Fatalf("expected 3 blocks, got %d", len(decoded.Blocks))
	}
	if decoded.Blocks[0].Type != "text" {
		t.Errorf("expected text block, got %s", decoded.Blocks[0].Type)
	}
	if decoded.Blocks[2].Children[0].TextProperties == nil {
		t.Error("child text properties should not be nil")
	}
	if decoded.Blocks[2].Children[0].Type != "text" {
		t.Errorf("expected child text block, got %s", decoded.Blocks[2].Children[0].Type)
	}
}
