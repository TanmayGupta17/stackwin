package models

import "time"

type Game struct {
	ID          string    `json:"id"`
	Player1ID   string    `json:"player1_id"`
	Player2ID   string    `json:"player2_id"`
	Player1Name string    `json:"player1_name"`
	Player2Name string    `json:"player2_name"`
	Board       [6][7]int `json:"board"` // 6 rows and 7 columns
	CurrentTurn string    `json:"current_turn"`
	Status      string    `json:"status"` // "active", "won", "draw"
	Winner      string    `json:"winner"` // ID of the winning player
	IsBot       bool      `json:"is_bot"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//Move struct
type Move struct {
	GameID    string    `json:"game_id"`
	PlayerID  string    `json:"player_id"`
	Column    int       `json:"column"`
	CreatedAt time.Time `json:"created_at"`
}

// GameResult stores completed games
type GameResult struct {
	GameID    string    `json:"game_id"`
	WinnerID  string    `json:"winner_id"`
	LoserID   string    `json:"loser_id"`
	Duration  int       `json:"duration"` // in seconds
	CreatedAt time.Time `json:"created_at"`
}
