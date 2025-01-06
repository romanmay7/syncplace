import React, { useState , useContext} from 'react';
import { useSelector, useDispatch } from "react-redux";
import { setChatMessageInStore } from "../CollabBoard/collabBoardSlice";
import { emitNewChatMessage } from "../wsocket/wsocketConn";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext
import {v4 as uuid} from "uuid";


function ChatMessage(msgId, roomId, timestamp, content, sender) {
    this.msgId = msgId;
    this.roomId = roomId;
    this.timestamp = timestamp;
    this.content = content;
    this.sender = sender;
  }


export const ChatComponent = ({setOpenedChatWindow }) => {

   //const[chat, setChat] = useState([]);
   const messages = useSelector(state =>state.collabBoard.chatMessages);
   const [message, setMessage] = useState("");
   const dispatch = useDispatch(); 
   const { currentRoom,userName } = useContext(AuthContext);

   const handleSubmit = (e) => {
      e.preventDefault();

      const timestamp = Date.now();
      // Create a new Date object from the timestamp
      const date = new Date(timestamp);
      // Using toISOString() for an ISO 8601 compliant format
      const formattedDate = date.toISOString();
    
      const chatMessage = new ChatMessage(
             uuid(),
             currentRoom,
             formattedDate, 
             message,
             userName
        );
      //Update Redux Store
      dispatch(setChatMessageInStore(chatMessage));
      //send message to all connected clients in the Room
      emitNewChatMessage(currentRoom,chatMessage);
      //setChat([... chat, message]);
      //setMessage("");
   }

   return (
    <div
      className="position-fixed top-0 h-100 text-white bg-dark"
      style={{ width:"350px", left:"%0 "}}
    >
      <button
        type="button"
        onClick={() => setOpenedChatWindow(false)}
        className="btn btn-light btn-block w-50 mt-5"
      >
        Close
      </button>
      <div className="w-100 mt-5 p-2  border border-1 border-white rounded-3"
           style={{height: "70%"}}
      >
        <div className="chat-messages">
        {messages.map((msg) => (
        <div key={msg.id} className="message">
          {/* Display message content, sender, timestamp, etc. */}
          <span className="sender">{msg.sender}: </span>
          <span className="content">{msg.content} </span>
          <span className="timestamp">{new Date(msg.timestamp).toLocaleString()}</span> 
        </div>
        ))}
        </div>

      </div>
      <form onSubmit={handleSubmit} className="w-100 mt-4  d-flex rounded-3">
        <input type ="text" 
        placeholder="Write Message" 
        className="w-100 border-0 rounded-0 py-2 px-4"
        style={{ width: "90%", }} 
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        />
        <button type="submit" className="btn btn-primary rounded-0">
            Send Message
        </button>
      </form>
    </div>
   );
};
