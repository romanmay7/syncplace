import { store } from "../store/store";
import { setAllBoardElementsInStore ,updateBoardElementInStore, setChatMessageInStore, setAllChatMessagesInStore } from "../CollabBoard/collabBoardSlice";


//import { v4 as uuidv4 } from 'uuid';

let wsocket;
const WEBSOCKET_URL ="http://localhost:3100";

export const Kind = {

    BOARD_STATE_UPDATE :"1",
    ELEMENT_UPDATE :"2",
    CONNECTED:"3",
    DISCONNCTED:"4",
    CHAT_MESSAGE:"5",
};

/*
let id_ = uuidv4();

const user = {
    id: id_,
    username: "user_" + id_,
};*/

//let roomId = "1";
//let roomName = "room_" + roomId ;

export const CreateNewRoom = async (roomId,roomName) => {
try{
    console.log("CreateNewRoom| roomId:"+roomId);
   //First CREATE ROOM for websocket communication (for group of clients)
    const res = await fetch(`${WEBSOCKET_URL}/ws/createRoom`, {
        method: "POST",
        headers: { "Content-Type": "application/json",},
        body: JSON.stringify({
            id: roomId,
            name: roomName,
        }),
    })

   if(res.ok) {
       console.log("The Room was created");  
       return true
    }
}
catch (err) {
    console.log(err)
    return false
}

return false

}

//CONNECT TO WEB SOCKET SERVER
export const connectToWSocketServer = (roomId,userName) => {

        /*const*/ wsocket = new WebSocket(`${WEBSOCKET_URL}/ws/joinRoom/${roomId}?username=${userName}`);
        if(wsocket.OPEN) {
            console.log("Connection to Web Socket has been established");
            
            //localStorage.setItem('current-room-id', roomId);
            //joinRoom(roomId);
        }

        wsocket.onmessage = function(event) {
            var messageData = event.data;
            const parsedData = JSON.parse(messageData);
        
            switch(parsedData.kind) {
                case(Kind.CONNECTED):
                console.log(messageData)
                break;
               case(Kind.ELEMENT_UPDATE):
                //console.log("ELEMENT_UPDATE");
                 store.dispatch(updateBoardElementInStore(parsedData.element))
                 break;
               case(Kind.BOARD_STATE_UPDATE):
                 {
                    console.log("BOARD_STATE_UPDATE");
                    if(parsedData.elements !== 'undefined' && parsedData.elements !== null)
                    {
                        store.dispatch(setAllBoardElementsInStore(parsedData.elements));
                    }
                    if(parsedData.chatMessages !== 'undefined' && parsedData.chatMessages !== null) 
                    {
                        store.dispatch(setAllChatMessagesInStore(parsedData.chatMessages));
                    }
                 }
                 break;
                case(Kind.CHAT_MESSAGE):
                 {
                    //console.log("CHAT_MESSAGE");
                    store.dispatch(setChatMessageInStore(parsedData.chatMessage))
                 }
                 break;
            }
        }

};


export const emitBoardElementUpdate = (roomId,elementData) => {
    if(wsocket.OPEN)
    {
        console.log("Connection to Web Socket is OPEN");
    }
    else
    {
        console.log("Connection to Web Socket is CLOSED");
    }
    wsocket.send(
        
        JSON.stringify({kind: Kind.ELEMENT_UPDATE, element: elementData,content: "Update Element",roomId})
    );
};

export const emitNewChatMessage = (roomId,message) => {
    if(wsocket.OPEN)
    {
        console.log("Connection to Web Socket is OPEN");
    }
    else
    {
        console.log("Connection to Web Socket is CLOSED");
    }
    console.log(message);
    wsocket.send(
        
        JSON.stringify({kind: Kind.CHAT_MESSAGE, chatMessage: message, content: "New Chat Message",roomId})
    );
};

