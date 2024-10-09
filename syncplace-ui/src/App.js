import React, {useEffect} from 'react';
import CollabBoard from './CollabBoard/CollabBoard';
import { connectToWSocketServer } from './wsocket/wsocketConn';



function App() {
useEffect(() => {
   connectToWSocketServer();
},[])

  return (
    <div >
     <CollabBoard />
    </div>
  );
}

export default App;
