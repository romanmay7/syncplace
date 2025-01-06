import { useState, useContext } from "react";
import { useNavigate, Link } from "react-router-dom";
import {connectToWSocketServer} from "../wsocket/wsocketConn";
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './index.css'; // Import your CSS file
import Logo from "../assets/SPLogo.png";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext
import {CreateNewRoom} from "../wsocket/wsocketConn";



const JoinRoom = ({ uuid }) => {
  const [roomId, setRoomId] = useState("");
  const { joinRoom , userName } = useContext(AuthContext); 

  const navigate = useNavigate();

  const handleJoinRoom = (e) => {
    e.preventDefault();
    
    CreateNewRoom(roomId,"").then(response => {
          if (response)
            {
    
                // Get UserName from Local Store
                //const userName = localStorage.getItem('syncplace-app-user');
                console.log("Room ID:"+roomId);
    
                connectToWSocketServer(roomId,userName);
                joinRoom(roomId);
            
                navigate("/collabboard");
            }
            else
            {
              console.log("There was a Problem to Create a channel for the  Room");
    
            }
          });
    // Get UserName from Local Store
    //const userName = localStorage.getItem('syncplace-app-user');
  

    navigate("/collabboard");
  };

  return (
    <div className="room-form-container"> {/* Container with styles */}
    <form className="form col-md-12 mt-5">
      <div className="form-group">
         <div className ="brand">
              <img src={Logo} alt="" />
         </div>
        <input
          type="text"
          className="form-control my-2"
          placeholder="Room Code"
          value={roomId}
          onChange={(e) => setRoomId(e.target.value)}
        />
      </div>
      <button
        type="submit"
        onClick={handleJoinRoom}
        className="mt-4 btn-primary btn-block form-control"
      >
        Join Room
      </button>
      <span >
        To Create another Room ? <Link to = "/createroom">Create Room</Link>
      </span>
    </form>
    </div>
  );
};

export default JoinRoom;
