import React, { useState , useContext} from 'react';
import { useSelector, useDispatch } from "react-redux";
import { setChatMessageInStore } from "../CollabBoard/collabBoardSlice";
import { emitNewChatMessage } from "../wsocket/wsocketConn";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext
import {v4 as uuid} from "uuid";
import { host } from "../utils/APIRoutes";
import './ChatComponent.css'


function ChatMessage(msgId, roomId, timestamp, content, sender, filePath) {
    this.msgId = msgId;
    this.roomId = roomId;
    this.timestamp = timestamp;
    this.content = content;
    this.sender = sender;
    this.filePath = filePath;
  }


export const ChatComponent = ({setOpenedChatWindow }) => {

   //const[chat, setChat] = useState([]);
   const messages = useSelector(state =>state.collabBoard.chatMessages);
   const [message, setMessage] = useState("");
   const [file, setFile] = useState(null);
   const dispatch = useDispatch(); 
   const { currentRoom,userName } = useContext(AuthContext);

//------------------------------------------------------------------------
   const handleSubmit = async (e) => {
      e.preventDefault();

      const timestamp = Date.now();
      // Create a new Date object from the timestamp
      const date = new Date(timestamp);
      // Using toISOString() for an ISO 8601 compliant format
      const formattedDate = date.toISOString();

      //Upload file to Server (file must be selected)
      let filePath = null;
      if (file) {
        const formData = new FormData();
        formData.append("file", file);

        const response = await fetch (host + '/api/upload', {
           method: 'POST',
           body: formData,
        });

        if(response.ok) {
          filePath = await response.text();
        }else {
          console.error("Error uploading file to Server: ",response.statusText);
         }
      }
    
      const chatMessage = new ChatMessage(
             uuid(),
             currentRoom,
             formattedDate, 
             message,
             userName,
             filePath
        );
      //Update Redux Store
      dispatch(setChatMessageInStore(chatMessage));
      //send message to all connected clients in the Room
      emitNewChatMessage(currentRoom,chatMessage);

      setMessage("");
      setFile(null);
   };
//--------------------------------------------------------------------------

const handleFileChange = (event) => {
  setFile(event.target.files[0]);
};

//--------------------------------------------------------------------------

   return (
    <div className="position-fixed top-0 h-100 text-dark bg-dark" style={{ width: "450px", left: "%0 " }}>
      <button type="button" onClick={() => setOpenedChatWindow(false)} className="btn btn-light btn-block w-50 mt-5">
        Close
      </button>
      <div className="w-100 mt-5 p-2 border border-1 border-white rounded-3" style={{ height: "70%" }}>
        <div className="chat-messages">
          {messages.map((msg) => (
            <div key={msg.id} className="message">
              <span className="sender">{msg.sender}: </span>
              <span className="content">
              <p>
                {msg.content}
                {msg.filePath && <a href={`${host}/uploads/${msg.filePath}`}>Download</a>}
              </p>
              </span>
              <span className="time-right">{new Date(msg.timestamp).toLocaleString()}</span>
            </div>
          ))}
        </div>
      </div>
      <form onSubmit={handleSubmit} className="w-100 mt-4 d-flex flex-column rounded-3"> 
        <input
          type="text"
          placeholder="Write Message"
          className="w-100 border-0 rounded-0 py-2 px-4 mb-2" 
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <div className="d-flex justify-content-between"> 
          <input type="file" id="fileInput" accept="*" onChange={handleFileChange} style={{ display: 'none' }} />
          <button type="button" className="btn btn-info rounded-0"
           onClick={() => document.getElementById('fileInput').click()}
           style= {{ marginLeft: '15px' }}>
            Attach File
          </button>
          <button type="submit" className="btn btn-primary rounded-0" style= {{ marginRight: '15px' }} >
            Send Message
          </button>
        </div>
      </form>
    </div>
   );
};
