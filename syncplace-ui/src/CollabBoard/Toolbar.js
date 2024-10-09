import React from "react";
import rectangleToolIcon from '../resources/icons/rectangle_tool.svg';
import { toolTypes } from "../definitions";
import { useDispatch, useSelector } from "react-redux";
import { setToolType } from "./collabBoardSlice";


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


const Toolbar = () => {
        return (
                <div className="toolbar_container">
                   <IconButton src={rectangleToolIcon} type ={toolTypes.RECTANGLE} /> 
                </div>
              );
   };

export default Toolbar;