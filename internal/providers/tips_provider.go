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
	log.Println("Getting tips", prediction, weather)
	data, err := json.Marshal(prediction)
	if err != nil {
		log.Println("Error marshalling prediction", err)
		return nil, err
	}
	log.Println("Tips data", string(data))
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
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body", err)
		return nil, err
	}
	log.Println("Tips response", string(body))
	var tip domain.Tips
	err = json.Unmarshal(body, &tip)
	if err != nil {
		log.Println("Error unmarshalling response body", err)
		return nil, err
	}
	log.Println("Tips datac", tip)
	return &tip, nil
}
