import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary

const HomeButton = () => {
    const navigate = useNavigate();
    const { handleRoomActions, currentRoom } = useWebSocketContext();

    const handleLeaveAndNavigate = () => {
        if (currentRoom) {
            handleRoomActions('leave', { roomName: currentRoom });
        }
        navigate('/'); // Navigate back to the lobby
    };

    return (
        <button onClick={handleLeaveAndNavigate}>Back To Lobby</button>
    );
};

export default HomeButton;
