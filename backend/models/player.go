package models

type Player struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Draws    int    `json:"draws"`
	ELO      int    `json:"elo_rating"`
}

type LeaderboardEntry struct {
	Players []Player `json:"players"`
}


