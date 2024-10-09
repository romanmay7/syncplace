import { store } from "../store/store";
import { setAllBoardElementsInStore ,updateBoardElementInStore } from "../CollabBoard/collabBoardSlice";
import { v4 as uuidv4 } from 'uuid';

let wsocket;
const WEBSOCKET_URL ="http://localhost:3100";

export const Kind = {

    BOARD_STATE_UPDATE :"1",
    ELEMENT_UPDATE :"2",
    CONNECTED:"3",
    DISCONNCTED:"4"
};

let id_ = uuidv4();

const user = {
    id: id_,
    username: "user_" + id_,
};

let roomId = "1";
let roomName = "room_" + roomId ;

export const CreateNewRoom = async () => {
try{
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
export const connectToWSocketServer = () => {

     // CreateNewRoom().then(response => {
     // if (response)
     // {   

        /*const*/ wsocket = new WebSocket(`${WEBSOCKET_URL}/ws/joinRoom/${roomId}?userId=${user.id}&username=${user.username}`);
        if(wsocket.OPEN) {
            console.log("Connection to Web Socket has been established");
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
                    //console.log("BOARD_STATE_UPDATE");
                    store.dispatch(setAllBoardElementsInStore(parsedData.elements))
                 }
                 break;
            }
        }
    //  }

   // })
};


export const emitBoardElementUpdate = (elementData) => {
    if(wsocket.OPEN)
    {
        console.log("Connection to Web Socket is OPEN");
    }
    else
    {
        console.log("Connection to Web Socket is CLOSED");
    }
    wsocket.send(
        
        JSON.stringify({kind: Kind.ELEMENT_UPDATE, element: elementData,content: "Update Element",roomId:"1"})
    );
};

