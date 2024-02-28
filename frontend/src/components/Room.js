import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary

const Room = () => {
    const location = useLocation();
    const roomName = location.pathname.split('/').pop();
    const { gameState, currentRoomPlayers, handleRoomActions, ws } = useWebSocketContext();
    const [numQuestions, setNumQuestions] = useState(5);
    const [joined, setJoined] = useState(false);


    useEffect(() => {
        const joinRoom = () => {
            // Make sure ws is not null before attempting to add event listeners or send messages
            if (ws) {
                if (ws.readyState === WebSocket.OPEN && !joined) {
                    handleRoomActions('join', { roomName });
                    setJoined(true)
                } else {
                }
            }
        };

        joinRoom();


    }, [ws, roomName, handleRoomActions, joined]);
    // This effect is for debugging purposes to log current room players
    useEffect(() => {
        console.log(currentRoomPlayers);
    }, [currentRoomPlayers]);

    const startGame = () => {
        handleRoomActions('startGame', { roomName, numQuestions });
    };

    const sendButtonIndex = (index) => {
        handleRoomActions('submitAnswer', { roomName, answerIdx: index });
    };

    return (
        <div>
            <h2>Room: {roomName}</h2>
            <p>Invite link: http://localhost:8000/room/{roomName}</p>
            {gameState.countdown && <h3>Game starts in: {gameState.countdown}</h3>}
            {currentRoomPlayers.length > 0 && (
                <>
                    <h3>Players in room:</h3>
                    <ul>
                        {currentRoomPlayers.map((player, index) => (
                            <li key={index}>{player}</li>
                        ))}
                    </ul>
                </>
            )}
            {!gameState.gameStarted && !gameState.gameOver && (
                <>
                    <label htmlFor="numQuestions">Number of Questions:</label>
                    <select
                        id="numQuestions"
                        value={numQuestions}
                        onChange={(e) => setNumQuestions(e.target.value)}
                    >
                        {[10, 15, 20].map((number) => (
                            <option key={number} value={number}>
                                {number}
                            </option>
                        ))}
                    </select>
                    <button onClick={startGame}>Start Game</button>
                </>
            )}
            {gameState.gameOver && (
                <div>
                    <h3>Game Over! Final Scores:</h3>
                    <ul>
                        {Object.entries(gameState.finalScores).map(([playerName, score], index) => (
                            <li key={index}>{playerName}: {score}</li>
                        ))}
                    </ul>
                    <select
                        id="numQuestions"
                        value={numQuestions}
                        onChange={(e) => setNumQuestions(e.target.value)}
                    >
                        {[10, 15, 20].map((number) => (
                            <option key={number} value={number}>
                                {number}
                            </option>
                        ))}
                    </select>
                    <button onClick={startGame}>Play Again</button>
                </div>
            )}
            {gameState.gameStarted && gameState.currentQuestion && (
                <div>
                    <h3>{gameState.currentQuestion.questionText}</h3>
                    {gameState.currentQuestion.options.map((option, index) => (
                        <button key={index} onClick={() => sendButtonIndex(index)}>
                            {option}
                        </button>
                    ))}
                </div>
            )}
            {!gameState.gameOver && gameState.answerFeedback && (
                <div>
                    <p>{`${gameState.answerFeedback.correct}`}</p>
                    <h4>Scores:</h4>
                    <ul>
                        {Object.entries(gameState.answerFeedback.scores).map(([playerName, score], index) => (
                            <li key={index}>{playerName}: {score}</li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default Room;
