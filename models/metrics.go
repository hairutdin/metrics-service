package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"` // "gauge" or "counter"
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
