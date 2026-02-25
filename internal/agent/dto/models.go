package dto

type RequestMetricUpdate struct {
	MainMetric
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MainMetric struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}
