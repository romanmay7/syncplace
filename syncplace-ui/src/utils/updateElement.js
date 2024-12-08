import { toolTypes } from "../definitions";
import { createBoardElement } from "./createElement";
import { store } from "../store/store";
import { setAllBoardElementsInStore } from "../CollabBoard/collabBoardSlice";
import { emitBoardElementUpdate } from "../wsocket/wsocketConn";


export const updateBoardElement = ({id, x1, x2, y1, y2, type, index }, elements, currentRoom) => {
    //create copy f elements array
    const elementsCopy = [...elements]

    switch(type) {
        case toolTypes.RECTANGLE:
            const updatedElement = createBoardElement({
                id,
                x1,
                y1,
                x2,
                y2,
                toolType: type,
            });
        
        //Replace the item with updated item
        elementsCopy[index] = updatedElement;

        //Update  all the elements in our Store
        store.dispatch(setAllBoardElementsInStore(elementsCopy));
        
        //Update all the Clients ,connected to the same Web Socket
        //const roomId = localStorage.getItem('current-room-id');  //Get roomID from localStore

        emitBoardElementUpdate(currentRoom,updatedElement);

        break;
        default:
            throw new Error('An Error occuried during elements update...')
    }
};