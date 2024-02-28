import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Lobby from './components/Lobby';
import Room from './components/Room';
import { WebSocketProvider } from './WebSocketContext'; // Import WebSocket context provider

function App() {
  const [currentPlayerName, setCurrentPlayerName] = useState('');
  const [playerName, setPlayerName] = useState('');

  const handleConnect = () => setPlayerName(currentPlayerName);

  if (!playerName) {
    return (
      <div className="App">
        <input
          type="text"
          placeholder="Enter your name"
          value={currentPlayerName}
          onChange={(e) => setCurrentPlayerName(e.target.value)}
        />
        <button onClick={handleConnect}>Connect</button>
      </div>
    );
  }

  return (
    <WebSocketProvider playerName={playerName}>
      <Router>
        <div className="App">
          <Routes>
            {/* The Lobby and Room components will access WebSocket functionalities through the context */}
            <Route path="/" element={<Lobby />} />
            <Route path="/room/:roomName" element={<Room />} />
          </Routes>
        </div>
      </Router>
    </WebSocketProvider>
  );
}

export default App;
