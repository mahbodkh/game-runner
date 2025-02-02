const ws = new WebSocket('ws://localhost:3000/ws');
let currentPlayer = null;
let currentSession = null;

ws.onopen = () => {
  console.log('WebSocket connection established');
  const playerName = prompt('Enter your name:');
  if (!playerName) {
    playerName = 'Guest' + Math.floor(Math.random() * 1000);
  }

  ws.send(playerName);
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received data:', data);
  switch (data.status) {
    case 'created':
      currentPlayer = data.player;
      currentSession = data.sessionId;
      updateSessionList();
      break;
    case 'started':
      updatePlayers(data.player1, data.player2);
      updateBoard(data.board);
      break;
    case 'available':
      displayAvailableSessions(data.sessions);
      break;
    case 'moved':
      updateBoard(data.board);
      break;
    case 'full':
      alert('Session is full! Please join another session.');
      break;
    case 'invalid':
      alert('Invalid session!');
      break;
    case 'completed':
      alert(
        `Game over! ${
          data.result === 'win' ? 'Winner: ' + data.winner : "It's a draw!"
        }`
      );
      break;
    case 'error':
      displayMessage(data.message);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket connection closed');
};

const cells = document.querySelectorAll('.cell');
cells.forEach((cell) => {
  cell.addEventListener('click', () => {
    if (currentSession && currentPlayer) {
      const index = cell.getAttribute('data-index');
      ws.send(
        JSON.stringify({
          action: 'move',
          index: index,
          sessionId: currentSession,
        })
      );
    }
  });
});

function updateBoard(board) {
  const cells = document.querySelectorAll('.cell');
  cells.forEach((cell, index) => {
    cell.textContent = board[index];
  });
}

function updatePlayers(player1, player2) {
  const playersDiv = document.getElementById('players');
  playersDiv.textContent = `${player1.name} (X) vs ${player2.name} (O)`;
}

function displayAvailableSessions(sessions) {
  const sessionList = document.getElementById('session-list');
  sessionList.innerHTML = '<h3>Available Sessions:</h3>';

  sessions.forEach((session) => {
    const sessionDiv = document.createElement('div');
    sessionDiv.className = 'session-item';
    sessionDiv.innerHTML = `
      <h4>Session ${session.sessionId}</h4>
      <p>Status: ${session.status}</p>
      <p>Players: ${session.players.join(', ')}</p>
      <button onclick="joinSession(${session.sessionId})">Join</button>
    `;
    sessionList.appendChild(sessionDiv);
  });
}

function joinSession(sessionId) {
  ws.send(
    JSON.stringify({
      action: 'join',
      sessionId: sessionId,
    })
  );
}

function updateSessionList() {
  ws.send(JSON.stringify({ action: 'getSessions' }));
}

function displayMessage(message) {
  const messagesDiv = document.getElementById('messages');
  const messageElement = document.createElement('p');
  messageElement.textContent = message;
  messagesDiv.appendChild(messageElement);
}
