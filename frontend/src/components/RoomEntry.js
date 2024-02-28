// RoomEntry.js
import React from 'react';

const RoomEntry = ({ roomName, onJoin }) => { // Removed setCurrentRoomPlayers from props since it's unused here
    const handleJoinClick = () => {
        onJoin(roomName); // Pass room name to join
    };

    return (
        <div>
            <span>{roomName}</span>
            <button onClick={handleJoinClick}>Join</button>
        </div>
    );
};

export default RoomEntry;
