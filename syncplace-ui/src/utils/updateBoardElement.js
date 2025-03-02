import { toolTypes } from "../definitions";
import { createBoardElement } from "./createBoardElement";
import { store } from "../store/store";
import { setAllBoardElementsInStore } from "../CollabBoard/collabBoardSlice";
import { emitBoardElementUpdate } from "../wsocket/wsocketConn";

//Utility function that updates the specific element on our Board
export const updateBoardElement = ({id, x1, x2, y1, y2, type,colour,fillMode,index }, elements, currentRoom) => {
    //create copy of elements array
    const elementsCopy = [...elements]

    switch(type) {
        case toolTypes.LINE :
        case toolTypes.CIRCLE :
        case toolTypes.RECTANGLE :
            const updatedElement = createBoardElement({
                id,
                x1,
                y1,
                x2,
                y2,
                colour:colour,
                fillMode:fillMode,
                toolType: type,
            });
        
        //Replace the item with updated item
        elementsCopy[index] = updatedElement;

        //Update  all the elements in our Store
        store.dispatch(setAllBoardElementsInStore(elementsCopy));
        
        //Update all the Clients ,connected to the same Web Socket, by sending them the updated element

        emitBoardElementUpdate(currentRoom,updatedElement);
        

       break;
        default:
            throw new Error('An Error occuried during elements update...')
    }
};