
import { useState, useEffect, useContext } from 'react';
import { AuthContext } from './AuthContext';
import { jwtDecode } from 'jwt-decode';

const AuthProvider = ({ children }) => {
  const [userName, setUserName] = useState('');
  const [token, setToken] = useState('');
  const [currentRoom, setCurrentRoom] = useState('');
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    if (storedToken && isTokenValid(storedToken)) {
      setToken(storedToken);
      setIsLoggedIn(true);
    } else {
        localStorage.removeItem('token'); //remove invalid token
    }
  }, []);

  const login = (newUserName, newToken) => {
    setUserName(newUserName);
    setToken(newToken);
    setIsLoggedIn(true);
    localStorage.setItem('token', newToken);
  };

  const joinRoom = (newCurrentRoom) => {
    setCurrentRoom(newCurrentRoom);
  };

  const leaveRoom = () => {
    setCurrentRoom('');
  };

  const logout = () => {
    setUserName('');
    setToken('');
    setIsLoggedIn(false);
    localStorage.removeItem('token');
  };

  const isTokenValid = (tokenToCheck) => {
    console.log ("Checking Token: " + tokenToCheck )
    if (!tokenToCheck) {
      return false;
    }
    try {
      const decodedToken = jwtDecode(tokenToCheck);
      const expirationDate = new Date(decodedToken.exp * 1000); // Multiply by 1000 to convert seconds to milliseconds
      console.log ("Decoded Token expiration: " + expirationDate )
      const currentTime = Date.now() / 1000;
      return decodedToken.exp > currentTime;
    } catch (error) {
      return false;
    }
  };

  const value = {
    userName,
    token,
    currentRoom,
    isLoggedIn,
    login,
    logout,
    joinRoom,
    leaveRoom,
    isTokenValid, // Expose isTokenValid
  };

  //AuthContext as a shared information container, value={value} is the actual information that is being placed inside the container.
  //{children} is all the components that will be allowed to look inside the container.
  
  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;