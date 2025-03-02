import React, {useContext, useState } from "react";
import rectangleToolIcon from '../resources/icons/rectangle_tool.svg';
import circleToolIcon from '../resources/icons/circle-svgrepo-com.svg';
import lineToolIcon from '../resources/icons/line-straight-svgrepo-com.svg';
import clearBoardIcon from '../resources/icons/Clear-Icon.png';
import saveIcon from '../resources/icons/save_512.png'
import { toolTypes } from "../definitions";
import { useDispatch, useSelector } from "react-redux";
import { setToolType, setColour, setFillModeInStore,setAllBoardElementsInStore } from "./collabBoardSlice";
import { host } from "../utils/APIRoutes";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext


//********************************************************************************************************************** */
const IconButton = ({src, type, isClearButton}) => {
 
    const dispatch = useDispatch();

    const selectedToolType = useSelector((state) => state.collabBoard.tool);

    const { currentRoom,userName } = useContext(AuthContext);

    const handleToolChange = () => {
       dispatch(setToolType(type));
    };

    const handleClearBoard = async() => {
      //Clear elements in Store
      dispatch(setAllBoardElementsInStore([]));
      
      //Clear Board State on Backend 
      try{
        const res = await fetch(`${host}/api/clearBoard/${currentRoom}/${userName}`, {
         method: "GET",
         headers: { "Content-Type": "application/json",},
         })
  
        if(res.ok) {
         console.log("The Room Board state was Cleared");
         return true
         }
       }
       catch (err) {
          console.log(err)
          return false
         }
    };
    

    return  (
    <button onClick={isClearButton ? handleClearBoard: handleToolChange} className={
        selectedToolType === type ? "toolbar_button_active" : "toolbar_button"
        }
    >
      <img width ='80%' height ='80%' src={src} />
    </button>
    );
};
//********************************************************************************************************************** */
const SaveButton =  ({src}) => {

   const { currentRoom } = useContext(AuthContext);

   const handleSaveBoard = async () => {
     //Save State
    try{
      const res = await fetch(`${host}/api/saveBoard/${currentRoom}`, {
       method: "GET",
       headers: { "Content-Type": "application/json",},
       })

      if(res.ok) {
       console.log("The Room Board state was Saved");
       alert("The Room Board state was Saved");  
       return true
       }
     }
     catch (err) {
        console.log(err)
        return false
       }
   };

   return  (
   <button onClick={handleSaveBoard} className={"toolbar_button"}
   >
     <img width ='80%' height ='80%' src={src} />
   </button>
   );
};

//********************************************************************************************************************** */
const Toolbar = () => {

        const dispatch = useDispatch();
        const [color, setColor] = useState('#000000'); // Add state for tool's color
        const [fillMode, setFillMode] = useState(true); // Add state to manage fill mode

        const handleFillModeChange = () => {
          setFillMode(!fillMode);
          console.log("Fill mode:", fillMode);
          // Dispatch setFillModeInStore event here
          dispatch(setFillModeInStore(fillMode));
        };

        const handleColorChange = (event) => {
          console.log("COLOUR:" + event.target.value);
          setColor(event.target.value); // Update color state first
          dispatch(setColour(event.target.value)); // Then dispatch action to update store 
        };

        return (
                <div className="toolbar_container">
                    <input 
                     type="color" 
                     value={color} 
                     onChange={handleColorChange} 
                    />
                   <IconButton src={lineToolIcon} type ={toolTypes.LINE} /> 
                   <IconButton src={rectangleToolIcon} type ={toolTypes.RECTANGLE} /> 
                   <IconButton src={circleToolIcon} type ={toolTypes.CIRCLE} />
                   <IconButton src={clearBoardIcon} isClearButton />
                   <SaveButton src={saveIcon} />
                   <input type="checkbox" checked={!fillMode} onChange={handleFillModeChange} />
                   <label>Fill Mode</label> 
                </div>
              );
   };

export default Toolbar;