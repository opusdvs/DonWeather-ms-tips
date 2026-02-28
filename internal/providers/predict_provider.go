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
	log.Println("Getting prediction", weather)
	data, err := json.Marshal(weather)
	if err != nil {
		return nil, err
	}
	log.Println("Prediction data", string(data))
	request, err := http.NewRequestWithContext(ctx, "POST", pp.apiUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := pp.clientPredict.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Println("Prediction response", string(body))
	var prediction domain.Prediction
	err = json.Unmarshal(body, &prediction)
	if err != nil {
		return nil, err
	}
	log.Println("Prediction", prediction)
	return &prediction, nil
}
