import styled from 'styled-components';

export const RoomContainer = styled.div`
    background-color: #f9f9f9;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    max-width: 600px;
    margin: 0 auto;
`;

export const Title = styled.h2`
    color: #333;
`;

export const Link = styled.p`
    color: #666;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    justify-content: center; /* Center horizontally */
    text-align: center; /* Center text within the container */
`;

export const CopyButton = styled.button`
    margin-left: 10px;
    padding: 5px 10px;
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

export const Countdown = styled.h3`
    color: #ff6347;
`;

export const PlayersList = styled.ul`
    list-style-type: none;
    padding: 0;
`;

export const PlayerItem = styled.li`
    color: #009688;
`;

export const Label = styled.label`
    color: #333;
    margin-right: 10px;
`;

export const Select = styled.select`
    padding: 8px;
    border-radius: 4px;
    border: 1px solid #ccc;
    margin-right: 10px;
`;

export const Button = styled.button`
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

export const Feedback = styled.div`
    margin-top: 20px;
`;

export const FeedbackText = styled.p`
    color: ${({ correct = 'default value' }) => (correct.includes('incorrect') ? '#f44336' : '#4caf50')};
`;

export const FeedbackScores = styled.ul`
    list-style-type: none;
    padding: 0;
`;

export const ScoreItem = styled.li`
    color: #009688;
`;

export const ButtonsContainer = styled.div`
    display: flex;
    flex-direction: column;
    margin-top: 10px;
`;

// -----

export const LobbyContainer = styled.div`
    background-color: #f9f9f9;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    max-width: 600px;
    margin: 0 auto;
`;

export const Form = styled.form`
    margin-bottom: 20px;
`;

export const Input = styled.input`
    padding: 10px;
    border-radius: 4px;
    border: 1px solid #ccc;
    margin-right: 10px;
`;

export const SubmitButton = styled.button`
    padding: 10px 20px;
    background-color: #4caf50;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    margin-left: 10px; /* Add margin here */

    &:hover {
        background-color: #45a049;
    }
`;

export const RoomList = styled.div`
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
`;
