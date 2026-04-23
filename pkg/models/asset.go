package models

type Asset struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string  `json:"type"`
	Children []Block `json:"children,omitempty"`

	// Type specific properties
	TextProperties  *TextProperties  `json:"textProperties,omitempty"`
	ImageProperties *ImageProperties `json:"imageProperties,omitempty"`
	TableProperties *TableProperties `json:"tableProperties,omitempty"`

	// Container specific
	Direction string  `json:"direction,omitempty"` // row or column
	Gap       float64 `json:"gap,omitempty"`
	Border    bool    `json:"border,omitempty"`

	// Common layout properties
	WidthPercent    *float64 `json:"widthPercent,omitempty"`
	MarginTop       *float64 `json:"marginTop,omitempty"`
	MarginBottom    *float64 `json:"marginBottom,omitempty"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
}

type TextProperties struct {
	Text            string  `json:"text"`
	FontFamily      string  `json:"fontFamily,omitempty"`
	FontSize        float64 `json:"fontSize"`
	FontWeight      string  `json:"fontWeight,omitempty"`
	FontColor       string  `json:"fontColor,omitempty"`
	BackgroundColor string  `json:"backgroundColor,omitempty"`
	LineHeight      float64 `json:"lineHeight,omitempty"`
	Align           string  `json:"align,omitempty"`
	MarginTop       float64 `json:"marginTop,omitempty"`
	MarginBottom    float64 `json:"marginBottom,omitempty"`
}

type ImageProperties struct {
	Path         string  `json:"path"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"` // 0 = auto
	MarginTop    float64 `json:"marginTop,omitempty"`
	MarginBottom float64 `json:"marginBottom,omitempty"`
	Align        string  `json:"align,omitempty"`
	OffsetX      float64 `json:"offsetX,omitempty"`
	OffsetY      float64 `json:"offsetY,omitempty"`
}

type TableProperties struct {
	Headers     []string   `json:"headers"`
	Rows        [][]string `json:"rows"`
	HeaderStyle *CellStyle `json:"headerStyle,omitempty"`
	RowStyle    *CellStyle `json:"rowStyle,omitempty"`
}

type CellStyle struct {
	CellHeight      float64 `json:"cellHeight,omitempty"`
	FontSize        float64 `json:"fontSize,omitempty"`
	FontWeight      string  `json:"fontWeight,omitempty"`
	FontColor       string  `json:"fontColor,omitempty"`
	BackgroundColor string  `json:"backgroundColor,omitempty"`
	Align           string  `json:"align,omitempty"`
}
