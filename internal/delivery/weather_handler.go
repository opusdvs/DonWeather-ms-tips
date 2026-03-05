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

	for _, output := range tips.Output {
		raw := extractJSON(output.Content[0].Text)

		// Пытаемся разобрать JSON вида {"tips":[{"text":"..."}, ...]}
		var parsed struct {
			Tips []struct {
				Text string `json:"text"`
			} `json:"tips"`
		}

		if err := json.Unmarshal([]byte(raw), &parsed); err != nil || len(parsed.Tips) == 0 {
			// fallback: возвращаем как есть, одной строкой
			tipsResponse.Tips = append(tipsResponse.Tips, domain.Tip{Text: raw})
			continue
		}

		for _, t := range parsed.Tips {
			if t.Text == "" {
				continue
			}
			tipsResponse.Tips = append(tipsResponse.Tips, domain.Tip{Text: t.Text})
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tipsResponse); err != nil {
		log.Println("Error encoding tips response", err)
	}
	log.Printf("tips response: %+v", tipsResponse)
}

func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end <= start {
		return text
	}
	return text[start : end+1]
}
