import React, { createContext, useContext, useState, useEffect } from 'react';

const API_BASE = process.env.REACT_APP_BACKEND_URL || "localhost:8080";
const WebSocketContext = createContext();

export const useWebSocketContext = () => useContext(WebSocketContext);

export const WebSocketProvider = ({ playerName, children }) => {
    const [ws, setWs] = useState(null);
    const [rooms, setRooms] = useState([]);
    const [currentRoom, setCurrentRoom] = useState('');
    const [currentRoomPlayers, setCurrentRoomPlayers] = useState([]);

    const [gameState, setGameState] = useState({
        gameStarted: false,
        gameOver: false,
        currentQuestion: null,
        answerFeedback: null,
        finalScores: {},
        countdown: null,
    });

    useEffect(() => {
        if (!playerName) return;

        const newWs = new WebSocket(`ws://${API_BASE}/ws`);

        newWs.onopen = () => {
            console.log('WebSocket Connected');
            newWs.send(JSON.stringify({ type: 'connect', playerName }));
        };

        newWs.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                switch (message.type) {
                    // case 'message':
                    //     if (message.data.startsWith("Joined room: ")) {
                    //         const joinedRoomName = message.data.replace("Joined room: ", "");
                    //         console.log(joinedRoomName)
                    //         if (joinedRoomName === currentRoom) {
                    //             // Call fetchRoomMetadata only if the joined room is the current room
                    //             fetchRoomMetadata(joinedRoomName);
                    //         }
                    //     }
                    //     break;
                    case 'playerListUpdate':
                        console.log("HURRRR")
                        console.log(message.data)
                        setCurrentRoomPlayers([...message.data]);

                        break;
                    case 'roomsList':
                        setRooms(message.data);
                        console.log(message.data)
                        break;
                    case 'roomUpdated':
                        if (message.roomName === currentRoom) {
                            setCurrentRoomPlayers(message.players);
                        }
                        break;
                    case 'question':
                        setGameState(prevState => ({
                            ...prevState,
                            gameStarted: true,
                            gameOver: false,
                            currentQuestion: message.data,
                            countdown: null
                        }));
                        break;
                    case 'answerFeedback':
                        setGameState(prevState => ({
                            ...prevState,
                            answerFeedback: message.data
                        }));
                        break;
                    case 'gameOver':
                        setGameState(prevState => ({
                            ...prevState,
                            gameOver: true,
                            finalScores: message.data.scores,
                            currentQuestion: null
                        }));
                        break;
                    case 'countdown':
                        setGameState(prevState => ({
                            ...prevState,
                            countdown: message.data
                        }));
                        break;
                    case 'error':
                        console.log(message.data)
                        break;
                    default:
                        console.log(message.data)
                        console.log('Unhandled message type:', message.type);
                }
            } catch (error) {
                console.error('Failed to parse message', error);
            }
        };

        setWs(newWs);

        return
    }, [playerName, currentRoom]);

    // Handling room actions
    const handleRoomActions = (action, data) => {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            console.error('WebSocket is not connected.');
            return;
        }
        let roomName = data.roomName
        if (roomName) {
            console.log(`room name ${roomName}`)
            setCurrentRoom(roomName);
        }
        // setCurrentRoom(roomName)
        // Send action with current room context if roomName is provided, else use currentRoom
        ws.send(JSON.stringify({ action, ...data }));

        switch (action) {
            case 'join':
                setCurrentRoom(roomName)
                break;
            case 'create':
                setCurrentRoom(roomName);
                break;
            case 'leave':
                setCurrentRoom('');
                setCurrentRoomPlayers([]);
                break;
            default:
                break;
        }
    };

    const contextValue = {
        ws,
        rooms,
        currentRoom,
        currentRoomPlayers,
        gameState,
        handleRoomActions,
    };

    return (
        <WebSocketContext.Provider value={contextValue}>
            {children}
        </WebSocketContext.Provider>
    );
};
