import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary
import HomeButton from "./HomeButton"
import { RoomContainer, Title, Link, CopyButton, Countdown, PlayersList, PlayerItem, Label, Select, Button, Feedback, FeedbackText, FeedbackScores, ScoreItem, ButtonsContainer } from './style';



const Room = () => {
    const location = useLocation();
    const roomName = location.pathname.split('/').pop();
    const { gameState, currentRoomPlayers, handleRoomActions, ws } = useWebSocketContext();
    const [numQuestions, setNumQuestions] = useState(10);
    const [joined, setJoined] = useState(false);

    const inviteLink = `http://localhost:3000/room/${roomName}`;

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
        const numQuestionsInt = parseInt(numQuestions);
        handleRoomActions('startGame', { roomName, numQuestions: numQuestionsInt });
        // handleRoomActions('startGame', { roomName, numQuestions });
    };

    const sendButtonIndex = (index) => {
        handleRoomActions('submitAnswer', { roomName, answerIdx: index });
    };

    const handleCopyLink = () => {
        navigator.clipboard.writeText(inviteLink)
            .then(() => alert('Invite link copied to clipboard'))
            .catch((error) => console.error('Failed to copy invite link: ', error));
    };

    return (
        <RoomContainer>
            <HomeButton />
            <Title>Room: {roomName}</Title>
            <Link>
                Invite link: {inviteLink}
                <CopyButton onClick={handleCopyLink}>Copy</CopyButton>
            </Link>
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
