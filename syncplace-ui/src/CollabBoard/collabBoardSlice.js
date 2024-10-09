
import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  tool:null,
  elements: []
};


export const collabBoardSlice = createSlice({
    name:'collabBoard',
    initialState,
    reducers: {
      setToolType: (state, action) => {
        state.tool = action.payload;
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
      },
    },
  });
  
export const  { setToolType, updateBoardElementInStore, setAllBoardElementsInStore } = collabBoardSlice.actions;
  
 // export default collabBoardSlice;