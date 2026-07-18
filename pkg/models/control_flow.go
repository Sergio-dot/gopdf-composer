// Package models defines the data structures for the PDF generation engine:
// document control flows, reusable assets, runtime contexts, and conditions.
package models

// ControlFlow describes the complete instruction set for generating a PDF
// document. It wraps a Document with its structure, headers, and footers.
type ControlFlow struct {
	Document Document `json:"document"`
}

// Document defines the structure and layout of the generated PDF, including
// sections of assets, header/footer assets, page size, orientation, and margins.
type Document struct {
	Structure    []Section        `json:"structure"`
	HeaderAssets []AssetReference `json:"headerAssets,omitempty"`
	FooterAssets []AssetReference `json:"footerAssets,omitempty"`

	PageSize     string  `json:"pageSize,omitempty"`
	Orientation  string  `json:"orientation,omitempty"`
	MarginLeft   float64 `json:"marginLeft,omitempty"`
	MarginTop    float64 `json:"marginTop,omitempty"`
	MarginRight  float64 `json:"marginRight,omitempty"`
	MarginBottom float64 `json:"marginBottom,omitempty"`
}

// Section is a group of asset references that are rendered together.
type Section struct {
	Assets []AssetReference `json:"assets"`
}

// AssetReference points to a reusable asset by ID and version, with optional
// conditions that determine whether the asset is included in the output.
type AssetReference struct {
	AssetID    string     `json:"assetId"`
	Version    string     `json:"version"`
	Conditions *Condition `json:"conditions,omitempty"`
}

// Condition defines a leaf or compound expression for conditional asset rendering.
//
// Leaf conditions use Field, Op, and Value (e.g., {"field": "age", "op": ">=", "value": 18}).
// Compound conditions use And, Or, or Not to combine sub-conditions.
// Supported leaf operators: ==, !=, >, <, >=, <=, in, contains.
type Condition struct {
	Field string `json:"field,omitempty"`
	Op    string `json:"op,omitempty"`
	Value any    `json:"value,omitempty"`

	And []Condition `json:"and,omitempty"`
	Or  []Condition `json:"or,omitempty"`
	Not *Condition  `json:"not,omitempty"`
}
