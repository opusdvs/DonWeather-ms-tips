package usecase

import (
	"context"

	"github.com/opusdvs/DonWeather-ms-tips/internal/domain"
)

type TipsService struct {
	predictProvider domain.PredictionProvider
	tipsProvider    domain.TipsProvider
}

func NewTipsService(
	predictProvider domain.PredictionProvider,
	tipsProvider domain.TipsProvider,
) *TipsService {
	return &TipsService{
		predictProvider: predictProvider,
		tipsProvider:    tipsProvider,
	}
}

func (ts *TipsService) GetTips(ctx context.Context, weather *domain.Weather) (*domain.Tip, error) {
	predict, err := ts.predictProvider.GetPrediction(ctx, weather)
	if err != nil {
		return nil, err
	}
	tip, err := ts.tipsProvider.GetTips(ctx, predict, weather)
	if err != nil {
		return nil, err
	}
	return tip, nil
}
