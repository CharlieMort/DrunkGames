import React, { useEffect, useState } from 'react';
import './App.css';
import RoomJoin from './comps/RoomJoin.tsx';
import { socket } from './comps/socket.tsx';
import { IClient, IPacket, IRoom } from './types.ts';
import Lobby from './comps/Lobby.tsx';
import Camera from './comps/Camera.jsx';

function App() {
  const [isConnected, setIsConnected] = useState<boolean>(true)
  const [packet, setPacket] = useState<IPacket>()
  const [client, setClient] = useState<IClient | undefined>(undefined)
  const [room, setRoom] = useState<IRoom | undefined>(undefined)
  const [imgUUID, setImgUUID] = useState<string>("")

  useEffect(() => {
    function onConnect() {
      console.log("CONNECTED")
      setIsConnected(true)
    }

    function onDisconnect() {
      console.log("DISCONNECTED")
      setIsConnected(false)
    }

    function onNewPacket(pkt) {
      console.log(pkt.data)
      setPacket(JSON.parse(pkt.data))
    }

    console.log("Setting Handlers")
    socket.onopen = onConnect
    socket.onclose = onDisconnect
    socket.onmessage = onNewPacket
  }, [])

  useEffect(() => {
    console.log(packet)
    if (packet !== undefined) {
      switch(packet.type) {
        case "toClient":
          setClient(JSON.parse(packet.data))
          console.log(JSON.parse(packet.data))
          break
        case "toRoom":
          setRoom(JSON.parse(packet.data))
          console.log(JSON.parse(packet.data))
          break;
      }
    }
  }, [packet])

  return (
    <div className="App">
      <h1>test</h1>
      {
        client !== undefined
        ? room === undefined 
            ? <RoomJoin client={client} />
            : <Lobby client={client} room={room} />
        : <h2>Waiting To Connect...</h2>
      }
    </div>
  );
}

export default App;