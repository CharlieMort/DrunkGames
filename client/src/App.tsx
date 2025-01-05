import React, { useEffect, useState } from 'react';
import './App.css';
import RoomJoin from './comps/RoomJoin.tsx';
import { Reconnect, socket } from './comps/socket.tsx';
import { IClient, IPacket, IRoom, ISettings, ISpyGame } from './types.ts';
import Lobby from './comps/Lobby.tsx';
import GameRouter from './comps/GameRouter.tsx';

function App() {
  const [isConnected, setIsConnected] = useState<boolean>(true)
  const [packet, setPacket] = useState<IPacket>()
  const [settings, setSettings] = useState<ISettings>({
    client: undefined,
    room: undefined,
    game: undefined,
  })
  const [x, sX] = useState(0)

  useEffect(() => {
    if (sessionStorage.tabID == undefined) {
      sessionStorage.tabID = crypto.randomUUID()
    }
    console.log(sessionStorage.tabID)

    function onConnect() {
      console.log("CONNECTED")
      socket.send(JSON.stringify({
        from: "0",
        to: "0",
        type: "setup",
        data: `clientconnect ${sessionStorage.tabID}`
      } as IPacket))
      setIsConnected(true)
    }

    function onDisconnect(e) {
      console.log("DISCONNECTED")
      console.log(e)
      setIsConnected(false)
      //Reconnect()
    }

    function onError(e) {
      console.log("ERROR")
      console.log(e)
    }

    function onNewPacket(pkt) {
      console.log(JSON.parse(pkt.data))
      setPacket(JSON.parse(pkt.data))
    }

    const refreshInt = setInterval(() => {
      sX(x+1)
      console.log("cum")
    }, 2000)

    console.log("Setting Handlers")
    socket.onopen = onConnect
    socket.onerror = onError
    socket.onclose = onDisconnect
    socket.onmessage = onNewPacket
    return () => clearInterval(refreshInt);
  }, [])

  useEffect(() => {
    if (packet === undefined || packet.type === "error") {
      return
    }
    let newSetting = {...settings}

    switch(packet.type) {
      case "toClient":
        console.log("toClient")
        newSetting.client = JSON.parse(packet.data)
        break
      case "toRoom":
        console.log("toRoom")
        newSetting.room = JSON.parse(packet.data)
        break;
      case "toGame":
        console.log("toGame")
        newSetting.game = JSON.parse(packet.data)
        break
    }
    console.log("NEWSETTING______________________________________________________________")
    console.log(newSetting)
    setSettings({...newSetting})
  }, [packet])

  return (
    <div className="App">
      <h1>🍺BigPint.com</h1>
      {
        settings.client !== undefined && isConnected
        ? settings.room === undefined 
          ? <RoomJoin client={settings.client} />
          : settings.game === undefined
            ? <Lobby client={settings.client} room={settings.room} />
            : <GameRouter room={settings.room} gameData={settings.game} />
        : <h2>Waiting To Connect...</h2>
      }
    </div>
  );
}

export default App;