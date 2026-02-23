package domain

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Condition   string  `json:"condition"`
}

type Tip struct {
	Text string `json:"text"`
}

type Prediction struct {
	PredictedActivity string  `json:"predicted_activity"`
	Confidence        float64 `json:"confidence"`
	ModelVersion      string  `json:"model_version"`
}
