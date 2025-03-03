import React, {useEffect} from 'react';
import CollabBoard from './CollabBoard/CollabBoard';
//import { connectToWSocketServer } from './wsocket/wsocketConn';
import {BrowserRouter as Router} from "react-router-dom"
import { Route, Routes} from "react-router";
import Registration from "./SignUpSignIn/Registration";
import Login from "./SignUpSignIn/Login";
import CreateRoom from './RoomsPortal/CreateRoom';
import JoinRoom from './RoomsPortal/JoinRoom';
import AuthProvider from './Auth/AuthProvider';
import ProtectedRoute from './Auth/ProtectedRoute';

function App() {
useEffect(() => {
  // connectToWSocketServer();
},[])
  return (
    <AuthProvider> {/* Wrap my entire app with AuthProvider to use AuthContext*/}
    <Router>
      <Routes>
         <Route path="/collabboard" element={ <ProtectedRoute><CollabBoard/></ProtectedRoute>} />
         <Route path="/signup" element={<Registration/>} />
         <Route path="/login" element={<Login/>} />
         <Route path="/createroom" element={<ProtectedRoute><CreateRoom/></ProtectedRoute>} />
         <Route path="/joinroom" element={<ProtectedRoute><JoinRoom/></ProtectedRoute>} />
      </Routes>
    </Router>
    </AuthProvider>
  );
}

export default App;
