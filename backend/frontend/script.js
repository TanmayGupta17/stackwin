let ws;
let gameState;
let playerID;
let currentGame;
let isPlayerOne = true;

function joinGame() {
    const username = document.getElementById('username').value;
    if (!username) {
        alert('Please enter username');
        return;
    }

    // Connect to WebSocket
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/ws`);

    ws.onopen = () => {
        ws.send(JSON.stringify({
            type: 'join',
            payload: { username: username }
        }));

        document.getElementById('waiting-msg').style.display = 'block';
    };

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        console.log('Received message:', msg);

        if (msg.type === 'game-state') {
            gameState = msg.payload.game;
            playerID = msg.payload.player_id;  // Get player ID from server!
            isPlayerOne = (playerID === gameState.player1_id);
            currentGame = msg.payload;

            console.log('Player ID:', playerID);
            console.log('Is Player One:', isPlayerOne);

            if (gameState.status === 'active') {
                document.getElementById('waiting-msg').style.display = 'none';
                showGameScreen();
                renderBoard();
                updateGameInfo();
            } else {
                showGameEndScreen();
            }
        }
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        alert('Connection error');
    };
}

function makeMove(column) {
    console.log('========== makeMove START ==========');
    console.log('Column:', column);
    console.log('gameState:', gameState);
    console.log('playerID:', playerID);
    console.log('WebSocket state:', ws ? ws.readyState : 'no ws', 'Open=', WebSocket.OPEN);

    if (!gameState) {
        console.error('ERROR: gameState is null!');
        alert('Game not initialized!');
        return;
    }

    if (!ws || ws.readyState !== WebSocket.OPEN) {
        console.error('ERROR: WebSocket not open');
        alert('Connection lost!');
        return;
    }

    if (gameState.current_turn !== playerID) {
        console.error('ERROR: Not your turn. Current turn:', gameState.current_turn, 'Player ID:', playerID);
        alert('Not your turn!');
        return;
    }

    // Find the row where the piece will land
    let landRow = -1;
    for (let row = 5; row >= 0; row--) {
        if (gameState.board[row][column] === 0) {
            landRow = row;
            break;
        }
    }

    // Check if column is full
    if (landRow === -1) {
        console.error('ERROR: Column is full!');
        alert('Column is full!');
        return;
    }

    console.log('Landing row:', landRow);

    // OPTIMISTIC UPDATE: Update the board immediately
    const playerPiece = isPlayerOne ? 1 : 2;  // Player 1 = 1, Player 2 = 2
    gameState.board[landRow][column] = playerPiece;
    console.log('Board updated at [' + landRow + '][' + column + '] with piece:', playerPiece);

    // Switch turn to the other player
    gameState.current_turn = isPlayerOne ? gameState.player2_id : gameState.player1_id;
    console.log('Turn switched to:', gameState.current_turn);

    renderBoard();
    updateGameInfo();

    // Send move to server
    const moveMsg = {
        type: 'move',
        payload: { column: column }
    };

    console.log('Sending move:', moveMsg);
    ws.send(JSON.stringify(moveMsg));
    console.log('Move sent!');
    console.log('========== makeMove END ==========');
}

function renderBoard() {
    const grid = document.getElementById('grid');
    grid.innerHTML = '';

    console.log('Rendering board, gameState:', gameState);

    for (let row = 0; row < 6; row++) {
        for (let col = 0; col < 7; col++) {
            const cell = document.createElement('div');
            cell.className = 'cell';
            cell.dataset.col = col;
            cell.dataset.row = row;

            const piece = gameState.board[row][col];
            console.log(`Cell [${row}][${col}] = ${piece}`);

            if (piece === 1) {
                cell.classList.add('player1');
                cell.textContent = '🟡';
            } else if (piece === 2) {
                cell.classList.add('player2');
                cell.textContent = '🔴';
            }

            // Add click handler to entire grid
            cell.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                console.log('Cell clicked at row:', row, 'col:', col);
                makeMove(col);
            });

            grid.appendChild(cell);
        }
    }
    console.log('Board render complete');
}

function updateGameInfo() {
    document.getElementById('player-name').textContent = gameState.player1_name;
    document.getElementById('current-turn').textContent =
        gameState.current_turn === gameState.player1_id ? 'Your Turn' : 'Opponent\'s Turn';

    if (gameState.status === 'won') {
        document.getElementById('game-status').textContent =
            gameState.winner === gameState.player1_id ? 'You Won! 🎉' : 'You Lost! 😢';
    } else if (gameState.status === 'draw') {
        document.getElementById('game-status').textContent = 'Draw!';
    }
}

function showGameScreen() {
    document.getElementById('login-screen').classList.remove('active');
    document.getElementById('game-screen').classList.add('active');
}

function showGameEndScreen() {
    document.getElementById('game-status').textContent =
        gameState.status === 'won' ? 'Game Ended' : 'Game Ended';
}

function leaveGame() {
    if (ws) {
        ws.send(JSON.stringify({ type: 'leave' }));
        ws.close();
    }
    location.reload();
}

function backToGame() {
    document.getElementById('leaderboard-screen').classList.remove('active');
    document.getElementById('game-screen').classList.add('active');
}
