import React, { useState } from 'react';
import RoomEntry from './RoomEntry'; // Make sure this path is correct for your project structure
import { useNavigate } from 'react-router-dom'; // Import useNavigate for programmatic navigation
import { useWebSocketContext } from '../WebSocketContext'; // Adjust the path as necessary

const Lobby = () => {
    const [newRoomName, setNewRoomName] = useState('');
    const navigate = useNavigate(); // Initialize useNavigate for navigation
    const { rooms, handleRoomActions } = useWebSocketContext();

    const handleCreateRoom = (e) => {
        e.preventDefault();
        if (newRoomName.trim()) {
            handleRoomActions('create', { roomName: newRoomName }); // Make sure to adjust the payload as per your backend's expectation
            navigate(`/room/${newRoomName}`);
            setNewRoomName('');
        }
    };

    const handleJoinRoom = (roomName) => {
        // Use handleRoomActions from context
        handleRoomActions('join', { roomName }); // Adjust the payload as necessary
        navigate(`/room/${roomName}`);
    };

    return (
        <div>
            <h2>Lobby</h2>
            <form onSubmit={handleCreateRoom}>
                <input
                    type="text"
                    placeholder="New room name"
                    value={newRoomName}
                    onChange={(e) => setNewRoomName(e.target.value)}
                />
                <button type="submit">Create Room</button>
            </form>
            <div>
                {rooms && rooms.map((room) => (
                    <RoomEntry
                        key={room.roomName}
                        roomName={room.roomName}
                        onJoin={() => handleJoinRoom(room.roomName)}
                    // Note: If setCurrentRoomPlayers isn't used in RoomEntry, you might not need to pass it
                    />
                ))}
            </div>
        </div>
    );
};

export default Lobby;
