import React, {useContext} from "react";
import chatIcon from '../resources/icons/vector-chat-icon-png_302635.jpg'
import { useDispatch, useSelector } from "react-redux";
//import { setToolType } from "./collabBoardSlice";
//import { host } from "../utils/APIRoutes";
import { AuthContext } from '../Auth/AuthContext'; // Import AuthContext


//********************************************************************************************************************** */

const ToggleChatButton =  ({src, setOpenedChatWindow }) => {

   //const { currentRoom } = useContext(AuthContext);


   return  (
   <button onClick={() => setOpenedChatWindow(true)} className={"toolbar_button"}
   >
     <img width ='80%' height ='80%' src={src} />
   </button>
   );
};

//********************************************************************************************************************** */
const FeaturesPanel = ({ setOpenedChatWindow }) => {
        return (
                <div className="features_panel">
                   <ToggleChatButton 
                   src={chatIcon} 
                   setOpenedChatWindow = {setOpenedChatWindow}
                   /> 
                </div>
              );
   };

export default FeaturesPanel;