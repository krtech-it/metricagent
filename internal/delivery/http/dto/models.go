package dto

type RequestUpdateMetric struct {
	MainMetric
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type RequestGetMetric struct {
	MainMetric
}

type ResponseGetMetric struct {
	MainMetric
	Value interface{} `json:"value"`
}

type MainMetric struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}
