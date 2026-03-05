package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
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

func (ts *TipsProvider) GetTips(ctx context.Context, prediction *domain.Prediction, weather *domain.Weather) (*domain.Tips, error) {
	// При минусовой температуре считаем вероятность дождя равной 0 (осадки — снег)
	predForLLM := *prediction
	if weather.Temperature < 0 {
		predForLLM.RainProbability = 0
	}
	// Передаём в LLM предсказание и текущую температуру (для правила «при минусовой — не зонт»)
	reqBody := struct {
		*domain.Prediction
		Temperature float64 `json:"temperature,omitempty"`
	}{Prediction: &predForLLM, Temperature: weather.Temperature}
	data, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("Error marshalling prediction", err)
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", ts.apiUrl, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error creating request", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.clientTips.Do(request)
	if err != nil {
		log.Println("Error doing request", err)
		return nil, err
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body", err)
		return nil, err
	}
	var tip domain.Tips
	err = json.Unmarshal(respBody, &tip)
	if err != nil {
		log.Println("Error unmarshalling response body", err)
		return nil, err
	}
	return &tip, nil
}
