package delivery

import (
	"net/http"

	"github.com/opusdvs/DonWeather-ms-tips/internal/usecase"
)

type WeatherHandler struct {
	tipsService *usecase.TipsService
}

func NewWeatherHandler(
	tipsService *usecase.TipsService,
) *WeatherHandler {
	return &WeatherHandler{
		tipsService: tipsService,
	}
}

func (wh *WeatherHandler) GetTips(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
