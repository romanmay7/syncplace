import { configureStore, getDefaultMiddleware } from '@reduxjs/toolkit';
import {collabBoardSlice} from  '../CollabBoard/collabBoardSlice';



export const store = configureStore({
  reducer: {
    collabBoard: collabBoardSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoreActions: ["collabBoard/setElements"],
        ignorePaths: ["collabBoard.elements"],
      },
    }),
});
 