package models

type ControlFlow struct {
	Document Document `json:"document"`
}

type Document struct {
	Structure    []Section        `json:"structure"`
	HeaderAssets []AssetReference `json:"headerAssets,omitempty"`
	FooterAssets []AssetReference `json:"footerAssets,omitempty"`
}

type Section struct {
	Assets []AssetReference `json:"assets"`
}

type AssetReference struct {
	AssetID    string     `json:"assetId"`
	Version    string     `json:"version"`
	Conditions *Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Field string `json:"field,omitempty"`
	Op    string `json:"op,omitempty"`
	Value any    `json:"value,omitempty"`

	And []Condition `json:"and,omitempty"`
	Or  []Condition `json:"or,omitempty"`
	Not *Condition  `json:"not,omitempty"`
}
