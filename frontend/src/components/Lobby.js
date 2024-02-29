import React, { useState } from 'react';
import RoomEntry from './RoomEntry'; // Make sure this path is correct for your project structure
import { useNavigate } from 'react-router-dom'; // Import useNavigate for programmatic navigation
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary
import { LobbyContainer, Title, Form, Input, SubmitButton, RoomList } from './style';

const Lobby = () => {
    const [newRoomName, setNewRoomName] = useState('');
    const navigate = useNavigate();
    const { rooms, handleRoomActions, resetGameState } = useWebSocketContext(); // Destructure resetGameState from context

    const handleCreateRoom = (e) => {
        e.preventDefault();
        if (newRoomName.trim()) {
            resetGameState(); // Reset game state before joining a new room
            handleRoomActions('create', { roomName: newRoomName });
            navigate(`/room/${newRoomName}`);
            setNewRoomName('');
        }
    };

    const handleJoinRoom = (roomName) => {
        resetGameState(); // Reset game state before joining a new room
        handleRoomActions('join', { roomName });
        navigate(`/room/${roomName}`);
    };

    return (
        <LobbyContainer>
            <Title>Lobby</Title>
            <Form onSubmit={handleCreateRoom}>
                <Input
                    type="text"
                    placeholder="New room name"
                    value={newRoomName}
                    onChange={(e) => setNewRoomName(e.target.value)}
                />
                <SubmitButton type="submit">Create Room</SubmitButton>
            </Form>
            <RoomList>
                {rooms && rooms.map((room) => (
                    <RoomEntry
                        key={room.roomName}
                        roomName={room.roomName}
                        onJoin={() => handleJoinRoom(room.roomName)}
                    />
                ))}
            </RoomList>
        </LobbyContainer>
    );
};

export default Lobby;

