import { createContext, useState, useEffect } from 'react';

export const AuthContext = createContext({
    userName: '',
    token: '',
    currentRoom: '',
    isLoggedIn: false,
    login: () => {},
    logout: () => {},
    joinRoom: () => {},
    leaveRoom: () => {},
});