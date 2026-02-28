package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/opusdvs/DonWeather-ms-tips/internal/domain"
	"github.com/opusdvs/DonWeather-ms-tips/internal/usecase"
)

type TipsHandler struct {
	tipsService *usecase.TipsService
}

func NewTipsHandler(
	tipsService *usecase.TipsService,
) *TipsHandler {
	return &TipsHandler{
		tipsService: tipsService,
	}
}

func (th *TipsHandler) GetTips(w http.ResponseWriter, r *http.Request) {
	var weatherRequest domain.WeatherRequest
	var tipsResponse domain.TipsResponse
	if err := json.NewDecoder(r.Body).Decode(&weatherRequest); err != nil {
		log.Println("Error decoding weather request", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	weather := &domain.Weather{
		Temperature: weatherRequest.Current.TempC,
		Humidity:    weatherRequest.Current.Humidity,
		WindSpeed:   weatherRequest.Current.WindKph,
		Pressure:    weatherRequest.Current.PressureMb,
	}

	tips, err := th.tipsService.GetTips(r.Context(), weather)
	if err != nil {
		log.Println("Error getting tips", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, tip := range tips.Output {
		tipsResponse.Tips = append(tipsResponse.Tips, domain.Tip{
			Text: extractJSON(tip.Content[0].Text),
		})
	}

	json.NewEncoder(w).Encode(tipsResponse)

	w.WriteHeader(http.StatusOK)
}

func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end <= start {
		return text
	}
	return text[start : end+1]
}
