
import { createSlice } from '@reduxjs/toolkit';
import { toolTypes } from "../definitions";

//Define Collaboration Board State that is managed by Redux Store
const initialState = {
  tool: toolTypes.LINE,
  colour:'#000000',
  fillMode: false,
  elements: [],
  chatMessages: [],
};

//********************************************************************************************************************** */
//Define functionality of Redux Store Management for Collaboration Board Component
//These functions can be called by using "dispatch" function that is part of reduxjs library
export const collabBoardSlice = createSlice({
    name:'collabBoard',
    initialState,
    reducers: {
      setToolType: (state, action) => {
        state.tool = action.payload;
      },
      setColour: (state, action) => {
        state.colour = action.payload;
      },
      setFillModeInStore: (state, action) => {
        state.fillMode = action.payload;
      },  
      updateBoardElementInStore: (state, action) => {
        console.log("updateBoardElementInStore");
        const {id} = action.payload;
        const elemIndex = state.elements.findIndex(element => element.id === id)
        //If not found
        if(elemIndex === -1) {
          state.elements.push(action.payload);
        }
        else {
          //If index has been found in elements array - Update element in our Array of Elements
          state.elements[elemIndex] = action.payload;
        }
      },  //Replacing the whole 'elements' array with new array from the payload
      setAllBoardElementsInStore: (state, action) => {
        state.elements = action.payload;
        console.log(state.elements);
      },//Set all  chat messages in the store
      setAllChatMessagesInStore: (state, action) => {
        var msgList = action.payload
        //sort chat messages by date
        msgList.sort(function(x, y){
          return x.timestamp - y.timestamp;
        })
        state.chatMessages = msgList;
        console.log(state.chatMessages);
      },
      //Add new Chat Message
      setChatMessageInStore: (state, action) => {
        state.chatMessages.push(action.payload);
        console.log(state.chatMessages);
      },
    },
  });
  
export const  { setToolType, setColour, setFillModeInStore, updateBoardElementInStore, setAllBoardElementsInStore, setChatMessageInStore, setAllChatMessagesInStore } = collabBoardSlice.actions;
  
 // export default collabBoardSlice;  