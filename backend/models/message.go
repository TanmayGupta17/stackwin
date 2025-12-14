package models

type Message struct {
	Type    string      `json:"type"` // "move", "join", "Leave", "game-state", "error"
	GameID  string      `json:"game_id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type JoinPayload struct {
	Username string `json:"username"`
}

type MovePayload struct {
	Column int `json:"column"`
}

type GameStatePayload struct {
	Game     *Game  `json:"game"`
	PlayerID string `json:"player_id"`
	Message  string `json:"message,omitempty"`
}
