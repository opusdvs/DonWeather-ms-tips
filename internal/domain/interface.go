package domain

import "context"

type PredictionProvider interface {
	GetPrediction(ctx context.Context, weather *Weather) (*Prediction, error)
}

type TipsProvider interface {
	GetTips(ctx context.Context, predoct *Prediction, weather *Weather) (*Tips, error)
}
