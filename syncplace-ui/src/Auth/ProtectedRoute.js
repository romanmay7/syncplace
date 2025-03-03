

import React, { useContext } from 'react';
import { Route, Navigate } from 'react-router-dom';
import { AuthContext } from './AuthContext';

//A Wrapper Component to ensure the protected access to his child components (used in App.js Router)
const ProtectedRoute = ({ children }) => {
  const { token, isLoggedIn, isTokenValid } = useContext(AuthContext);
  
  //In case user is not logged in or his token has expired we will redirect him to login page
  if (!isLoggedIn || !isTokenValid(token)) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

export default ProtectedRoute;