import { useState, useEffect, useContext } from 'react';
import { AuthContext } from './AuthContext';

const AuthProvider = ({ children }) => {
    const [userName, setUserName] = useState('');
    const [token, setToken] = useState('');
    const [currentRoom, setCurrentRoom] = useState('');
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    useEffect(() => {
        // Check for existing token or other authentication methods
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            setToken(storedToken);
            setIsLoggedIn(true);
            // Fetch user details and current app if needed
        }
    }, []);

    const login = (newUserName, newToken) => {
        setUserName(newUserName);
        setToken(newToken);

        setIsLoggedIn(true);
        localStorage.setItem('token', newToken); // Store token for future use
    };

    const joinRoom =(newCurrentRoom) => {
        setCurrentRoom(newCurrentRoom);
    };

    const leaveRoom =() => {
        setCurrentRoom('');
    };

    const logout = () => {
        setUserName('');
        setToken('');
        setIsLoggedIn(false);
        localStorage.removeItem('token');
    };

    const value = {
        userName,
        token,
        currentRoom,
        isLoggedIn,
        login,
        logout,
        joinRoom,
        leaveRoom

    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthProvider;