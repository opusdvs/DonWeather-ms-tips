package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/opusdvs/DonWeather-ms-tips/internal/delivery"
	"github.com/opusdvs/DonWeather-ms-tips/internal/providers"
	"github.com/opusdvs/DonWeather-ms-tips/internal/usecase"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	predictApiUrl := os.Getenv("PREDICT_API_URL")
	if predictApiUrl == "" {
		log.Fatal("PREDICT_API_URL is not set")
	}
	llmApiUrl := os.Getenv("LLM_API_URL")
	if llmApiUrl == "" {
		log.Fatal("LLM_API_URL is not set.")
	}

	predictProvider := providers.NewPredictionProvider(predictApiUrl)
	tipsProvider := providers.NewTipsProvider(llmApiUrl)
	tipsService := usecase.NewTipsService(predictProvider, tipsProvider)
	tipsHandler := delivery.NewTipsHandler(tipsService)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/v1/tips", tipsHandler.GetTips)

	mainMux := http.NewServeMux()
	mainMux.Handle("/", apiMux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      apiMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func(server *http.Server) {
		log.Println("Starting API server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start API server: %v", err)
		}
	}(server)

	<-appCtx.Done()
	log.Println("Shutting down API server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown API server: %v", err)
	}
	log.Println("API server shutdown complete")
}
