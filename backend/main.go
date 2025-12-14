package main

import (
	"log"
	"net/http"

	"4-in-a-row/config"
	"4-in-a-row/handlers"
	"4-in-a-row/services"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	// Initialize services
	gameService := services.NewGameService()
	botService := services.NewBotService(gameService)
	matchmakingService := services.NewMatchmakingService(cfg.MatchmakingTimeout)

	// Initialize Kafka (disabled for now - causing delays)
	var analyticsService *services.AnalyticsService
	// producer := kafka.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	// defer producer.Close()
	// analyticsService := services.NewAnalyticsService(producer)
	
	// For now, create a nil-safe analytics service
	analyticsService = services.NewAnalyticsService(nil)

	// Initialize handler
	gameHandler := handlers.NewGameHandler(gameService, botService, matchmakingService, analyticsService)

	// Set up routes
	router := mux.NewRouter()

	// WebSocket endpoint
	router.HandleFunc("/ws", gameHandler.HandleWebSocket)

	// Static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

	log.Printf("Server starting on %s\n", cfg.Port)
	http.ListenAndServe(cfg.Port, router)
}
