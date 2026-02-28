package usecase

import (
	"context"
	"log"

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

func (ts *TipsService) GetTips(ctx context.Context, weather *domain.Weather) (*domain.Tips, error) {
	predict, err := ts.predictProvider.GetPrediction(ctx, weather)
	if err != nil {
		log.Println("Error getting prediction", err)
		return nil, err
	}
	tips, err := ts.tipsProvider.GetTips(ctx, predict, weather)
	if err != nil {
		log.Println("Error getting tips", err)
		return nil, err
	}
	return tips, nil
}
