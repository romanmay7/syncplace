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

function App() {
useEffect(() => {
  // connectToWSocketServer();
},[])
  return (
    <AuthProvider> {/* Wrap my entire app with AuthProvider to use AuthContext*/}
    <Router>
      <Routes>
         <Route path="/collabboard" element={<CollabBoard/>} />
         <Route path="/signup" element={<Registration/>} />
         <Route path="/login" element={<Login/>} />
         <Route path="/createroom" element={<CreateRoom/>} />
         <Route path="/joinroom" element={<JoinRoom/>} />
      </Routes>
    </Router>
    </AuthProvider>
  );
}

export default App;
