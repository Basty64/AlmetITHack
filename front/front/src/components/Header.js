import React from 'react';
import { Link } from 'react-router-dom';

function Header() {
    return (
        <header className="header">
            <div className="container">
                <h1>Online Shop</h1>
                <nav>
                <ul>
                    <li><a href='/' className="styled-link">Home</a></li>
                    <li><a href='/login' className="styled-link">Login</a></li>
                    <li><a href='/register' className="styled-link">Register</a></li>
                </ul>
            </nav>
            </div>
        </header>
    );
}

export default Header;
