import React, {useContext} from "react";
import rectangleToolIcon from '../resources/icons/rectangle_tool.svg';
import saveIcon from '../resources/icons/save_512.png'
import { toolTypes } from "../definitions";
import { useDispatch, useSelector } from "react-redux";
import { setToolType } from "./collabBoardSlice";
import { host } from "../utils/APIRoutes";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext


const IconButton = ({src, type}) => {
 
    const dispatch = useDispatch();

    const selectedToolType = useSelector((state) => state.collabBoard.tool);

    const handleToolChange = () => {
       dispatch(setToolType(type));
    };

    return  (
    <button onClick={handleToolChange} className={
        selectedToolType === type ? "toolbar_button_active" : "toolbar_button"
        }
    >
      <img width ='80%' height ='80%' src={src} />
    </button>
    );
};

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


const Toolbar = () => {
        return (
                <div className="toolbar_container">
                   <IconButton src={rectangleToolIcon} type ={toolTypes.RECTANGLE} /> 
                   <SaveButton src={saveIcon} /> 
                </div>
              );
   };

export default Toolbar;