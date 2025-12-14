# 4-in-a-Row Game

A real-time multiplayer Connect Four game built with **Go** backend and **vanilla JavaScript** frontend. Play against other players or challenge the AI bot!

## Features

✨ **Multiplayer Gameplay** - Play against another person in real-time via WebSocket  
🤖 **AI Bot** - Challenge an intelligent bot that automatically joins if no opponent is available  
⚡ **Smooth UI** - Optimistic updates for instant feedback on moves  
🔄 **Real-time Sync** - Board state synchronized across all connected players  
📱 **Responsive Design** - Works on desktop and mobile browsers  
🎮 **Turn-based Logic** - Clear turn indicators and game state management  

## Tech Stack

**Backend:**
- Go 1.21+
- Gorilla WebSocket
- PostgreSQL (database ready)
- Docker & Docker Compose

**Frontend:**
- HTML5
- CSS3
- Vanilla JavaScript (no frameworks)

## Project Structure

```
4-in-a-row/
├── backend/
│   ├── main.go                 # Server entry point
│   ├── go.mod
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── handlers/
│   │   └── websocket_handlers.go # WebSocket logic
│   ├── services/
│   │   ├── game_service.go     # Game rules & logic
│   │   ├── bot_service.go      # AI bot implementation
│   │   ├── matchmaking_service.go # Player pairing
│   │   └── analytics_service.go   # Event logging
│   ├── models/
│   │   ├── game.go
│   │   ├── message.go
│   │   └── player.go
│   ├── database/
│   │   └── db.go
│   └── kafka/
│       ├── producer.go
│       └── consumer.go
├── frontend/
│   ├── index.html
│   ├── style.css
│   └── script.js
├── docker-compose.yml
└── README.md
```

## Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Docker & Docker Compose** - [Download](https://www.docker.com/products/docker-desktop)
- **Git** - [Download](https://git-scm.com/)

## Setup Instructions

### Quick Start (Local Development)

1. **Clone repository:**
   ```bash
   git clone https://github.com/TanmayGupta17/stackwin.git
   cd stackwin
   ```

2. **Start services:**
   ```bash
   docker-compose up -d
   ```

3. **Run backend:**
   ```bash
   cd backend
   go run main.go
   ```

4. **Open game:**
   Navigate to http://localhost:8080

### Production Deployment

**Deploy to Railway in 5 minutes!** 🚀

See [HOSTING.md](HOSTING.md) for complete deployment guide.

Or click the button below to deploy directly:

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template?code=TanmayGupta17/stackwin)

### Detailed Setup (Local Development)

## How to Play

### Single Player (vs Bot)

1. Open http://localhost:8080
2. Enter your username
3. Click "Join Game"
4. Wait 10 seconds for the bot to join automatically
5. Start playing!

### Multiplayer (vs Another Player)

1. **Tab 1**: Open http://localhost:8080 → Enter username → Click "Join Game"
2. **Tab 2**: Open http://localhost:8080 → Enter a **different username** → Click "Join Game" (within 10 seconds)
3. Both tabs should now show the game board
4. Player 1 (yellow 🟡) goes first
5. Take turns clicking columns to drop pieces

### Game Rules

- Drop pieces into columns from the top
- Pieces fall to the lowest available row
- First player to get 4 pieces in a row (horizontal, vertical, or diagonal) wins
- If the board fills with no winner, it's a draw

## Configuration

Create a `.env` file in the project root to customize settings:

```env
PORT=:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/four_in_a_row?sslmode=disable
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=game-events
```

## Troubleshooting

### Port Already in Use

If port 8080 is already in use, change it in `.env`:
```env
PORT=:3000
```

### Docker Services Not Starting

```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart

# Full cleanup and restart
docker-compose down -v
docker-compose up -d
```

### WebSocket Connection Lost

- Check that the backend server is running
- Hard refresh browser (Ctrl+Shift+R)
- Check browser console for errors (F12)

### Build Errors

```bash
# Download dependencies
go mod tidy

# Clear Go cache
go clean -cache

# Try building again
go build
```

## API Endpoints

### WebSocket Endpoint
- `ws://localhost:8080/ws` - Real-time game communication

### Message Types

**Join Game**
```json
{
  "type": "join",
  "payload": { "username": "PlayerName" }
}
```

**Make Move**
```json
{
  "type": "move",
  "payload": { "column": 3 }
}
```

**Game State Update** (from server)
```json
{
  "type": "game-state",
  "game_id": "uuid",
  "payload": {
    "game": { ... },
    "player_id": "PlayerID",
    "message": "Move accepted"
  }
}
```

## Development

### Running Tests

```bash
cd backend
go test ./...
```

### Code Format

```bash
go fmt ./...
```

### Build Binary

```bash
cd backend
go build -o 4-in-a-row
./4-in-a-row
```

## Features in Development

- ✅ Real-time multiplayer
- ✅ Bot AI
- ⏳ Leaderboard system
- ⏳ Game history
- ⏳ User authentication
- ⏳ Elo rating system
- ⏳ Mobile app

## Performance Optimizations

- **Optimistic UI Updates** - Moves appear instantly without waiting for server
- **Asynchronous Operations** - Bot moves and logging run in background
- **WebSocket Ping** - Keeps connections alive with 30-second heartbeat
- **Efficient Message Routing** - Only broadcasts to connected players

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

Having issues? 
- Check the Troubleshooting section above
- Open an issue on GitHub
- Check server logs: `docker-compose logs -f`

## Author

Created with ❤️ by [Your Name]

---

**Happy Playing! 🎮**
