package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Create tables
	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	playersTable := `
	CREATE TABLE IF NOT EXISTS players (
		id VARCHAR(100) PRIMARY KEY,
		username VARCHAR(100) UNIQUE,
		wins INT DEFAULT 0,
		losses INT DEFAULT 0,
		draws INT DEFAULT 0,
		elo INT DEFAULT 1000,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	gamesTable := `
	CREATE TABLE IF NOT EXISTS games (
		id VARCHAR(100) PRIMARY KEY,
		player1_id VARCHAR(100),
		player2_id VARCHAR(100),
		winner_id VARCHAR(100),
		status VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (player1_id) REFERENCES players(id),
		FOREIGN KEY (player2_id) REFERENCES players(id)
	);`

	db.Exec(playersTable)
	db.Exec(gamesTable)
}
