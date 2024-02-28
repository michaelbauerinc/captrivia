import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary
import styled from 'styled-components';

const RoomContainer = styled.div`
    background-color: #f9f9f9;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    max-width: 600px;
    margin: 0 auto;
`;

const Title = styled.h2`
    color: #333;
`;

const Link = styled.p`
    color: #666;
    margin-bottom: 20px;
`;

const Countdown = styled.h3`
    color: #ff6347;
`;

const PlayersList = styled.ul`
    list-style-type: none;
    padding: 0;
`;

const PlayerItem = styled.li`
    color: #009688;
`;

const Label = styled.label`
    color: #333;
    margin-right: 10px;
`;

const Select = styled.select`
    padding: 8px;
    border-radius: 4px;
    border: 1px solid #ccc;
    margin-right: 10px;
`;

const Button = styled.button`
    padding: 10px 20px;
    background-color: #4caf50;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    margin-top: 10px;

    &:hover {
        background-color: #45a049;
    }
`;

const Feedback = styled.div`
    margin-top: 20px;
`;

const FeedbackText = styled.p`
    color: ${({ correct }) => (correct.includes('incorrect') ? '#f44336' : '#4caf50')};
`;

const FeedbackScores = styled.ul`
    list-style-type: none;
    padding: 0;
`;

const ScoreItem = styled.li`
    color: #009688;
`;

const ButtonsContainer = styled.div`
    display: flex;
    flex-direction: column;
    margin-top: 10px;
`;

const Room = () => {
    const location = useLocation();
    const roomName = location.pathname.split('/').pop();
    const { gameState, currentRoomPlayers, handleRoomActions, ws } = useWebSocketContext();
    const [numQuestions, setNumQuestions] = useState(5);
    const [joined, setJoined] = useState(false);

    useEffect(() => {
        const joinRoom = () => {
            if (ws) {
                if (ws.readyState === WebSocket.OPEN && !joined) {
                    handleRoomActions('join', { roomName });
                    setJoined(true);
                }
            }
        };

        joinRoom();
    }, [ws, roomName, handleRoomActions, joined]);

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
        <RoomContainer>
            <Title>Room: {roomName}</Title>
            <Link>Invite link: http://localhost:8000/room/{roomName}</Link>
            {gameState.countdown && <Countdown>Game starts in: {gameState.countdown}</Countdown>}
            {currentRoomPlayers.length > 0 && (
                <>
                    <h3>Players in room:</h3>
                    <PlayersList>
                        {currentRoomPlayers.map((player, index) => (
                            <PlayerItem key={index}>{player}</PlayerItem>
                        ))}
                    </PlayersList>
                </>
            )}
            {!gameState.gameStarted && !gameState.gameOver && (
                <>
                    <Label htmlFor="numQuestions">Number of Questions:</Label>
                    <Select
                        id="numQuestions"
                        value={numQuestions}
                        onChange={(e) => setNumQuestions(e.target.value)}
                    >
                        {[10, 15, 20].map((number) => (
                            <option key={number} value={number}>
                                {number}
                            </option>
                        ))}
                    </Select>
                    <Button onClick={startGame}>Start Game</Button>
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
                    <Select
                        id="numQuestions"
                        value={numQuestions}
                        onChange={(e) => setNumQuestions(e.target.value)}
                    >
                        {[10, 15, 20].map((number) => (
                            <option key={number} value={number}>
                                {number}
                            </option>
                        ))}
                    </Select>
                    <Button onClick={startGame}>Play Again</Button>
                </div>
            )}
            {gameState.gameStarted && gameState.currentQuestion && (
                <div>
                    <h3>{gameState.currentQuestion.questionText}</h3>
                    <ButtonsContainer>
                        {gameState.currentQuestion.options.map((option, index) => (
                            <Button key={index} onClick={() => sendButtonIndex(index)}>
                                {option}
                            </Button>
                        ))}
                    </ButtonsContainer>
                </div>
            )}
            {!gameState.gameOver && gameState.answerFeedback && (
                <Feedback>
                    <FeedbackText correct={gameState.answerFeedback.correct}>
                        {gameState.answerFeedback.correct}
                    </FeedbackText>
                    <h4>Scores:</h4>
                    <FeedbackScores>
                        {Object.entries(gameState.answerFeedback.scores).map(([playerName, score], index) => (
                            <ScoreItem key={index}>{playerName}: {score}</ScoreItem>
                        ))}
                    </FeedbackScores>
                </Feedback>
            )}
        </RoomContainer>
    );
};

export default Room;
