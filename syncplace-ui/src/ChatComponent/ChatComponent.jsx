import React from 'react';

export const ChatComponent = ({setOpenedChatWindow }) => {
   return (
    <div
      className="position-fixed top-0 h-100 text-white bg-dark"
      style={{ width:"350px", left:"%0 "}}
    >
      <button
        type="button"
        onClick={() => setOpenedChatWindow(false)}
        className="btn btn-light btn-block w-50 mt-5"
      >
        Close
      </button>
      <div className="w-100 mt-5 pt-5">
           Chat Component
      </div>

    </div>
   );
};
