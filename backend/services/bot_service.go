package services

import (
	"4-in-a-row/models"
)

type BotService struct {
	gameservice *GameService
}

func NewBotService(gs *GameService) *BotService {
	return &BotService{
		gameservice: gs,
	}
}

func (bs *BotService) MakeBotMove(game *models.Game) int {
	// Quick win check
	botpiece := 2
	for col := 0; col < 7; col++ {
		if canPlaceInColumn(&game.Board, col) {
			row := getLowestRow(&game.Board, col)
			game.Board[row][col] = botpiece

			if bs.gameservice.checkWin(game.Board, botpiece) {
				game.Board[row][col] = 0
				return col
			}
			game.Board[row][col] = 0
		}
	}

	// Quick block check - only check critical positions
	playerPiece := 1
	for col := 0; col < 7; col++ {
		if canPlaceInColumn(&game.Board, col) {
			row := getLowestRow(&game.Board, col)
			game.Board[row][col] = playerPiece

			if bs.gameservice.checkWin(game.Board, playerPiece) {
				game.Board[row][col] = 0
				return col
			}
			game.Board[row][col] = 0
		}
	}

	// Prefer center columns - fastest heuristic
	preferredOrder := []int{3, 4, 2, 5, 1, 6, 0}
	for _, col := range preferredOrder {
		if canPlaceInColumn(&game.Board, col) {
			return col
		}
	}

	return -1
}

func canPlaceInColumn(board *[6][7]int, column int) bool {
	if column < 0 || column >= 7 {
		return false
	}
	return board[0][column] == 0
}

func getLowestRow(board *[6][7]int, column int) int {
	for row := 5; row >= 0; row-- {
		if board[row][column] == 0 {
			return row
		}
	}
	return -1 // Column is full
}
