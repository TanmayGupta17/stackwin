package services

import (
	"time"

	"4-in-a-row/kafka"
)

type AnalyticsService struct {
	producer *kafka.KafkaProducer
}

func NewAnalyticsService(producer *kafka.KafkaProducer) *AnalyticsService {
	return &AnalyticsService{producer: producer}
}

func (as *AnalyticsService) LogMove(gameID, playerID string, column int) {
	if as.producer == nil {
		return
	}
	event := kafka.GameEvent{
		EventType: "move",
		GameID:    gameID,
		PlayerID:  playerID,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"column": column,
		},
	}
	as.producer.SendEventAsync(event)
}

func (as *AnalyticsService) LogGameEnd(gameID, winnerID, status string) {
	if as.producer == nil {
		return
	}
	event := kafka.GameEvent{
		EventType: "game_end",
		GameID:    gameID,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"winner": winnerID,
			"status": status,
		},
	}
	as.producer.SendEventAsync(event)
}

func (as *AnalyticsService) LogGameAbandoned(gameID, playerID string) {
	if as.producer == nil {
		return
	}
	event := kafka.GameEvent{
		EventType: "game_abandoned",
		GameID:    gameID,
		PlayerID:  playerID,
		Timestamp: time.Now().Unix(),
	}
	as.producer.SendEventAsync(event)
}
