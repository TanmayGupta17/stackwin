package services

import (
	"4-in-a-row/models"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrGameNotFound   = errors.New("game not found")
	ErrGameNotActive  = errors.New("game is not active")
	ErrNotPlayersTurn = errors.New("not player's turn")
	ErrInvalidColumn  = errors.New("invalid column")
	ErrColumnFull     = errors.New("column is full")
)

type GameService struct {
	games map[string]*models.Game
	mu    sync.RWMutex
}

func NewGameService() *GameService {
	return &GameService{
		games: make(map[string]*models.Game),
	}
}

func (gs *GameService) CreateGame(player1ID, player1Name, player2ID, player2Name string, isBot bool) *models.Game {
	game := &models.Game{
		ID:          uuid.New().String(),
		Player1ID:   player1ID,
		Player1Name: player1Name,
		Player2ID:   player2ID,
		Player2Name: player2Name,
		CurrentTurn: player1ID,
		Status:      "active",
		IsBot:       isBot,
	}
	gs.mu.Lock()
	gs.games[game.ID] = game
	gs.mu.Unlock()
	return game
}

func (gs *GameService) StoreGame(game *models.Game) {
	gs.mu.Lock()
	gs.games[game.ID] = game
	gs.mu.Unlock()
}

func (gs *GameService) MakeMove(gameID, playerID string, column int) (*models.Game, error) {
	gs.mu.Lock()
	game, exists := gs.games[gameID]
	gs.mu.Unlock()
	if !exists {
		return nil, ErrGameNotFound
	}
	if game.Status != "active" {
		return nil, ErrGameNotActive
	}
	if game.CurrentTurn != playerID {
		return nil, ErrNotPlayersTurn
	}
	if column < 0 || column >= 7 {
		return nil, ErrInvalidColumn
	}

	// Find the lowest empty row in the specified column
	row := -1
	for r := 5; r >= 0; r-- {
		if game.Board[r][column] == 0 {
			row = r
			break
		}
	}

	if row == -1 {
		return nil, ErrColumnFull
	}

	//place piece
	piece := 1
	if playerID == game.Player2ID {
		piece = 2
	}
	game.Board[row][column] = piece

	if gs.checkWin(game.Board, piece) {
		game.Status = "won"
		game.Winner = playerID
	} else if gs.isBoardFull(game.Board) {
		game.Status = "draw"
	} else {
		if game.CurrentTurn == game.Player1ID {
			game.CurrentTurn = game.Player2ID
		} else {
			game.CurrentTurn = game.Player1ID
		}
	}
	gs.mu.Lock()
	gs.games[game.ID] = game
	gs.mu.Unlock()
	return game, nil
}

func (gs *GameService) checkWin(board [6][7]int, piece int) bool {
	// Check horizontal, vertical, diagonal
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			if board[row][col] == piece {
				if gs.checkDirection(board, row, col, piece, 0, 1) || // horizontal
					gs.checkDirection(board, row, col, piece, 1, 0) || // vertical
					gs.checkDirection(board, row, col, piece, 1, 1) || // diagonal right
					gs.checkDirection(board, row, col, piece, 1, -1) { // diagonal left
					return true
				}
			}
		}
	}
	return false
}

func (gs *GameService) checkDirection(board [6][7]int, row, col int, piece int, dRow, dCol int) bool {
	// count := 0
	for i := 0; i < 4; i++ {
		r := row + i*dRow
		c := col + i*dCol
		if r < 0 || r >= 6 || c < 0 || c >= 7 || board[r][c] != piece {
			return false
		}
	}
	return true
}

func (gs *GameService) isBoardFull(board [6][7]int) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

func (gs *GameService) GetGame(gameID string) (*models.Game, error) {
	gs.mu.RLock()
	game, exists := gs.games[gameID]
	gs.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("game not found")
	}
	return game, nil
}

func (gs *GameService) DeleteGame(gameID string) {
	gs.mu.Lock()
	delete(gs.games, gameID)
	gs.mu.Unlock()
}
