package services

import (
	"4-in-a-row/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type WaitingPlayer struct {
	ID        string
	Name      string
	Timestamp time.Time
	Channel   chan *models.Game
}

type MatchmakingService struct {
	WaitingPlayers map[string]*WaitingPlayer
	mu             sync.RWMutex
	timeout        int // in seconds
}

func NewMatchmakingService(timeout int) *MatchmakingService {
	ms := &MatchmakingService{
		WaitingPlayers: make(map[string]*WaitingPlayer),
		timeout:        timeout,
	}
	go ms.cleanupExpiredPlayers()
	return ms
}

func (ms *MatchmakingService) AddPlayer(playerID, playerName string) *models.Game {
	ms.mu.Lock()
	var match *WaitingPlayer
	for _, wp := range ms.WaitingPlayers {
		match = wp
		break
	}
	if match != nil {
		delete(ms.WaitingPlayers, match.ID)
		ms.mu.Unlock()

		game := &models.Game{
			ID:          uuid.New().String(),
			Player1ID:   match.ID,
			Player1Name: match.Name,
			Player2ID:   playerID,
			Player2Name: playerName,
			CurrentTurn: match.ID,
			Status:      "active",
			IsBot:       false,
		}
		match.Channel <- game
		return game
	}

	wp := &WaitingPlayer{
		ID:        playerID,
		Name:      playerName,
		Timestamp: time.Now(),
		Channel:   make(chan *models.Game),
	}
	ms.WaitingPlayers[playerID] = wp
	ms.mu.Unlock()

	select {
	case game := <-wp.Channel:
		return game

	case <-time.After(time.Duration(ms.timeout) * time.Second):
		ms.mu.Lock()
		delete(ms.WaitingPlayers, playerID)
		ms.mu.Unlock()
		game := &models.Game{
			ID:          uuid.New().String(),
			Player1ID:   playerID,
			Player1Name: playerName,
			Player2ID:   "bot",
			Player2Name: "Bot",
			CurrentTurn: playerID,
			Status:      "active",
			IsBot:       true,
		}
		return game
	}
}

func (ms *MatchmakingService) cleanupExpiredPlayers() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		ms.mu.Lock()
		now := time.Now()
		for id, wp := range ms.WaitingPlayers {
			if now.Sub(wp.Timestamp) > time.Duration(ms.timeout)*time.Second {
				delete(ms.WaitingPlayers, id)
			}
		}
		ms.mu.Unlock()
	}
}

func (ms *MatchmakingService) RemovePlayer(playerID string) {
	ms.mu.Lock()
	delete(ms.WaitingPlayers, playerID)
	ms.mu.Unlock()
}