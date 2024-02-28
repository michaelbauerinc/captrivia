import React from 'react';
import { Link } from 'react-router-dom';

const HomeButton = () => {
    return (
        <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <button>Back To Lobby</button>
        </Link>
    );
};

export default HomeButton;
