import React, {useRef, useLayoutEffect, useState, useContext} from "react";
import { useSelector, useDispatch } from "react-redux";
import Toolbar from "./Toolbar";
import FeaturesPanel from "./FeaturesPanel";
import rough from 'roughjs/bundled/rough.esm';
import { actions, toolTypes } from "../definitions";
import { createBoardElement, updateBoardElement, drawBoardElement } from "../utils";
import {v4 as uuid} from "uuid";
import { updateBoardElementInStore } from "./collabBoardSlice";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext
import { ChatComponent } from '../ChatComponent/ChatComponent';

let selectedBoardElement;


//Used during element update ,when we are moving the mouse, resizing our element ("handleMouseMove" function)
const setSelectedBoardElement = (el) => {
    selectedBoardElement = el;
};

//********************************************************************************************************************** */

//Define Collaboration Board Component
const CollabBoard = () => {
    const canvasRef = useRef();
  
    const { currentRoom } = useContext(AuthContext);

    //Getting objects from Store's State
    const toolType = useSelector(state =>state.collabBoard.tool);
    const colour = useSelector(state =>state.collabBoard.colour);
    const fillMode = useSelector(state =>state.collabBoard.fillMode);
    const elements = useSelector(state =>state.collabBoard.elements);

    const [action, setAction] = useState(null);
    const[openedChatWindow, setOpenedChatWindow] = useState(false);
    const dispatch = useDispatch(); 
    //---------------------------------------------------------------------------------
    //The RENDERING Part of Collaboration Board Canvas goes here
    //useLayoutEffect is a version of useEffect that fires before the browser repaints the screen.
    useLayoutEffect(() => {
        //GET CANVAS
        const canvas = canvasRef.current;
        //GET CONTEXT
        const ctx = canvas.getContext("2d");

        //CLEAR CANVAS STATE
        ctx.clearRect(0, 0, canvas.width, canvas.height)

        const rc = rough.canvas(canvas);

        //RENDER ALL ELEMENTS ON CANVAS
        elements.forEach(element => {
            drawBoardElement({roughCanvas: rc, context: ctx, element });
        });

        //TEST RENDERING
        //rc.rectangle(15, 16, 222, 200);
        //rc.rectangle(20, 20, 333, 600);

    }, [elements]); //Everytime 'elements' array changes - useLayoutEffect() will run
   //-----------------------------------------------------------------------------------

   //On Mouse Click ,Crating new Board Element on Canvas,based on mouse coordinates
    const handleMouseDown = (event) => {
         const {clientX, clientY} = event;
         
         console.log(toolType);

         if(toolType === toolTypes.RECTANGLE || toolType === toolTypes.CIRCLE || toolType === toolTypes.LINE ){
           setAction(actions.DRAWING);
        }

         console.log(clientX, clientY);
         
         //Create new element by providing mouse coordinates on Canvas , selected tool type, colour and fillMode parameters,
         //assignng it new ID by using "uuid()" function
         const element = createBoardElement({
            x1:clientX,
            y1:clientY,
            x2:clientX,
            y2:clientY,
            toolType,
            colour,
            fillMode,
            id: uuid(),
         });
         
         //reference the newly created element as the "selectedBoardElement"
         setSelectedBoardElement(element);
         //Update our Redux Store with newly created element
         dispatch(updateBoardElementInStore(element));
         console.log(element);
    };
   //------------------------------------------------------------
    const handleMouseUp = () => {
           setAction(null)
           setSelectedBoardElement(null);

        }
    //------------------------------------------------------------
    const handleMouseMove = (event) => {
        const {clientX, clientY} = event;

        if (action === actions.DRAWING) {
            //Find index of the Selected Board Element in "elements" array
            const index = elements.findIndex((el) =>el.id == selectedBoardElement.id)
            //If Found
            if(index !== -1) {
                updateBoardElement({
                    index,
                    id:elements[index].id,
                    x1: elements[index].x1, //OLD x
                    y1: elements[index].y1, //OLD y
                    x2: clientX, //NEW x
                    y2: clientY, //NEW y
                    type: elements[index].type,
                    colour: colour,
                    fillMode:fillMode,
                }, 
                elements,
                currentRoom
              );
            }
        }
    }
    //------------------------------------------------------------

    return (
     <>
       <Toolbar />
       <FeaturesPanel
         setOpenedChatWindow = {setOpenedChatWindow}
       />
       {openedChatWindow && <ChatComponent setOpenedChatWindow = {setOpenedChatWindow}/>}
       <canvas
       onMouseDown={handleMouseDown}
       onMouseUp={handleMouseUp}
       onMouseMove ={handleMouseMove}
         ref={canvasRef} 
         width={window.innerWidth}
         height={window.innerHeight}
       />
     </>
    );
};

export default CollabBoard;