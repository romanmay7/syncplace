import React, {useEffect} from 'react';
import CollabBoard from './CollabBoard/CollabBoard';
import { connectToWSocketServer } from './wsocket/wsocketConn';
import {BrowserRouter as Router} from "react-router-dom"
import { Route, Routes} from "react-router";
import Registration from "./SignUpSignIn/Registration";
import Login from "./SignUpSignIn/Login";

function App() {
useEffect(() => {
   connectToWSocketServer();
},[])
  return (
    <Router>
      <Routes>
         <Route path="/collabboard" element={<CollabBoard/>} />
         <Route path="/signup" element={<Registration/>} />
         <Route path="/login" element={<Login/>} />
      </Routes>
    </Router>
  );
}

export default App;
