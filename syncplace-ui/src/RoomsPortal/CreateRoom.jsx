import { useState, useContext } from "react";
import { useNavigate, Link } from "react-router-dom";
import {CreateNewRoom} from "../wsocket/wsocketConn";
import {connectToWSocketServer} from "../wsocket/wsocketConn";
import {v4 as uuid} from "uuid";
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import './index.css'; // Import your CSS file
import Logo from "../assets/SPLogo.png";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext

const CreateRoom = () => {
  const [roomId, setRoomId] = useState(uuid());
  const [name, setRoomName] = useState("");
  const { joinRoom,userName } = useContext(AuthContext); 

  const navigate = useNavigate();

  const handleCreateRoom = (e) => {
    e.preventDefault();

     CreateNewRoom(roomId,name).then(response => {
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
          console.log("There was a Problem to Create a New Room");

        }
      });

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
            placeholder="Room Name"
            value={name}
            onChange={(e) => setRoomName(e.target.value)}
          />
        </div>
        <div className="form-group border">
          <div className="input-group d-flex align-items-center justify-content-center">
            <input
              type="text"
              value={roomId}
              className="form-control my-2 border-0"
              disabled
              placeholder="Room Code"
            />
            <div className="input-group-append">
              <button
                className="btn btn-primary btn-sm me-1"
                onClick={() => setRoomId(uuid())}
                type="button"
              >
                Generate
              </button>
              <button
                className="btn btn-outline-danger btn-sm me-2"
                type="button"
              >
                Copy
              </button>
            </div>
          </div>
        </div>
        <button
          type="submit"
          onClick={handleCreateRoom}
          
          className="mt-4 btn-primary btn-block form-control"
        >
          Create New Room
        </button>
        <span >
                To Join another Room ? <Link to = "/joinroom">Join Room</Link>
        </span>
      </form>
    </div>
  );

};

export default CreateRoom;
