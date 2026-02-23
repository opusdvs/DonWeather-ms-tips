package providers

import (
	"context"
	"net/http"
	"time"

	"github.com/opusdvs/DonWeather-ms-tips/internal/domain"
)

type TipsProvider struct {
	apiUrl     string
	clientTips *http.Client
}

func NewTipsProvider(apiUrl string) *TipsProvider {
	return &TipsProvider{
		apiUrl: apiUrl,
		clientTips: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ts *TipsProvider) GetTips(ctx context.Context, predoct *domain.Prediction, weather *domain.Weather) (*domain.Tip, error) {
	return &domain.Tip{}, nil
}
