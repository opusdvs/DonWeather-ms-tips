package providers

import (
	"context"
	"net/http"
	"time"

	"github.com/opusdvs/DonWeather-ms-tips/internal/domain"
)

type PredictionProvider struct {
	apiUrl        string
	clientPredict *http.Client
}

func NewPredictionProvider(apiUrl string) *PredictionProvider {
	return &PredictionProvider{
		apiUrl: apiUrl,
		clientPredict: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (pp *PredictionProvider) GetPrediction(ctx context.Context, weather *domain.Weather) (*domain.Prediction, error) {
	return &domain.Prediction{}, nil
}
