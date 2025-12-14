package handlers

import (
	"4-in-a-row/models"
	"4-in-a-row/services"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GameHandler struct {
	gameService      *services.GameService
	botService       *services.BotService
	matchService     *services.MatchmakingService
	analyticsService *services.AnalyticsService
	clients          map[string]*websocket.Conn
	mu               sync.RWMutex
}

func NewGameHandler(gs *services.GameService, bs *services.BotService, ms *services.MatchmakingService, ans *services.AnalyticsService) *GameHandler {
	return &GameHandler{
		gameService:      gs,
		botService:       bs,
		matchService:     ms,
		analyticsService: ans,
		clients:          make(map[string]*websocket.Conn),
	}
}

func (gh *GameHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	var playerID string
	var gameID string

	// Ping ticker to keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Done channel to signal goroutine to stop
	done := make(chan bool)
	defer close(done)

	// Start ping goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Printf("Ping error: %v\n", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error for player %s: %v\n", playerID, err)
			break
		}

		switch msg.Type {
		case "join":
			payload := msg.Payload.(map[string]interface{})
			username := payload["username"].(string)
			playerID = username + "_" + generateID()

			log.Printf("Player joining: %s (playerID: %s)\n", username, playerID)

			// Register player IMMEDIATELY and ALWAYS
			gh.mu.Lock()
			gh.clients[playerID] = conn
			gh.mu.Unlock()
			log.Printf("Stored connection for playerID: %s\n", playerID)

			// Matchmaking
			game := gh.matchService.AddPlayer(playerID, username)
			gameID = game.ID

			log.Printf("Game created: Player1ID=%s, Player2ID=%s, IsBot=%v\n", game.Player1ID, game.Player2ID, game.IsBot)

			// Store the game in game service
			gh.gameService.StoreGame(game)

			// Send game state
			message := "Game started! Your turn."
			if game.IsBot {
				message = "Playing against Bot. Your turn!"
			}

			// Broadcast with ONLY the sender's connection to confirm receipt
			senderResponse := models.Message{
				Type:    "game-state",
				GameID:  game.ID,
				Payload: models.GameStatePayload{Game: game, PlayerID: playerID, Message: message},
			}
			conn.WriteJSON(senderResponse)

			// Broadcast to other player with delay to ensure they're ready
			time.Sleep(100 * time.Millisecond)
			gh.broadcastToOthers(game, playerID, message)

		case "move":
			if gameID == "" {
				log.Println("Error: gameID is empty")
				gh.sendError(conn, "game not started")
				continue
			}

			log.Println("Move message received. Payload:", msg.Payload)
			payload := msg.Payload.(map[string]interface{})
			column := int(payload["column"].(float64))
			log.Printf("Move: gameID=%s, playerID=%s, column=%d\n", gameID, playerID, column)

			game, err := gh.gameService.MakeMove(gameID, playerID, column)
			if err != nil {
				log.Printf("MakeMove error: %v\n", err)
				gh.sendError(conn, err.Error())
				continue
			}

			log.Printf("Move successful. Game status: %s, Current turn: %s\n", game.Status, game.CurrentTurn)
			// Broadcast updated state to both players
			gh.broadcastGameState(game, "Move accepted")

			// Log event asynchronously to avoid blocking
			go gh.analyticsService.LogMove(gameID, playerID, column)

			// If bot's turn, make bot move asynchronously
			if game.IsBot && game.Status == "active" && game.CurrentTurn == "bot" {
				go func() {
					time.Sleep(1 * time.Second) // Delay for better UX
					log.Println("Making bot move...")
					botCol := gh.botService.MakeBotMove(game)
					log.Printf("Bot column: %d\n", botCol)
					if botCol >= 0 {
						updatedGame, err := gh.gameService.MakeMove(gameID, "bot", botCol)
						if err == nil {
							gh.broadcastGameState(updatedGame, "Bot moved")
							go gh.analyticsService.LogMove(gameID, "bot", botCol)

							// Check if game is finished
							if updatedGame.Status != "active" {
								gh.analyticsService.LogGameEnd(gameID, updatedGame.Winner, updatedGame.Status)
							}
						}
					}
				}()
			} else if game.Status != "active" {
				// Check if game is finished for human vs human
				go gh.analyticsService.LogGameEnd(gameID, game.Winner, game.Status)
			}

		case "leave":
			if gameID != "" {
				gh.analyticsService.LogGameAbandoned(gameID, playerID)
				gh.gameService.DeleteGame(gameID)
			}
			break
		}
	}

	gh.mu.Lock()
	delete(gh.clients, playerID)
	gh.mu.Unlock()
	log.Printf("Player disconnected: %s\n", playerID)
}

func (gh *GameHandler) broadcastGameState(game *models.Game, message string) {
	gh.mu.RLock()
	conn1 := gh.clients[game.Player1ID]
	conn2 := gh.clients[game.Player2ID]
	gh.mu.RUnlock()

	if conn1 != nil {
		response := models.Message{
			Type:    "game-state",
			GameID:  game.ID,
			Payload: models.GameStatePayload{Game: game, PlayerID: game.Player1ID, Message: message},
		}
		if err := conn1.WriteJSON(response); err != nil {
			log.Printf("Error writing to Player1: %v\n", err)
		}
	}
	if conn2 != nil && game.Player2ID != "bot" {
		response := models.Message{
			Type:    "game-state",
			GameID:  game.ID,
			Payload: models.GameStatePayload{Game: game, PlayerID: game.Player2ID, Message: message},
		}
		if err := conn2.WriteJSON(response); err != nil {
			log.Printf("Error writing to Player2: %v\n", err)
		}
	}
}

func (gh *GameHandler) broadcastToOthers(game *models.Game, senderID string, message string) {
	gh.mu.RLock()
	var otherID string
	var otherConn *websocket.Conn

	if senderID == game.Player1ID && game.Player2ID != "bot" {
		otherID = game.Player2ID
		otherConn = gh.clients[game.Player2ID]
	} else if senderID == game.Player2ID {
		otherID = game.Player1ID
		otherConn = gh.clients[game.Player1ID]
	}
	gh.mu.RUnlock()

	if otherConn != nil {
		response := models.Message{
			Type:    "game-state",
			GameID:  game.ID,
			Payload: models.GameStatePayload{Game: game, PlayerID: otherID, Message: message},
		}
		if err := otherConn.WriteJSON(response); err != nil {
			log.Printf("Error writing to other player: %v\n", err)
		}
	} else {
		log.Printf("Other player connection not found. Sender: %s, Game: P1=%s P2=%s\n", senderID, game.Player1ID, game.Player2ID)
	}
}

func (gh *GameHandler) sendError(conn *websocket.Conn, errMsg string) {
	response := models.Message{
		Type:    "error",
		Payload: map[string]string{"error": errMsg},
	}
	conn.WriteJSON(response)
}

func generateID() string {
	return uuid.NewString()
}
