import React from 'react';
import styled from 'styled-components';

const RoomEntryContainer = styled.div`
    background-color: #ffffff;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px; /* Add margin bottom */
    display: flex;
    justify-content: space-between; /* Add space between room name and join button */
    align-items: center; /* Center align items */
`;

const RoomName = styled.span`
    font-size: 1.2rem;
`;

const JoinButton = styled.button`
    padding: 10px 20px;
    background-color: #4caf50;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s ease;

    &:hover {
        background-color: #45a049;
    }
`;

const RoomEntry = ({ roomName, onJoin }) => {
    const handleJoinClick = () => {
        onJoin(roomName); // Pass room name to join
    };

    return (
        <RoomEntryContainer>
            <RoomName>{roomName}</RoomName>
            <JoinButton onClick={handleJoinClick}>Join</JoinButton>
        </RoomEntryContainer>
    );
};

export default RoomEntry;
