import React, { useEffect, useState } from 'react';
import './App.css';
import RoomJoin from './comps/RoomJoin.tsx';
import { socket } from './comps/socket.tsx';
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

  useEffect(() => {
    function onConnect() {
      console.log("CONNECTED")
      if (sessionStorage.tabID == undefined) {
        sessionStorage.tabID = crypto.randomUUID()
      }
      console.log(sessionStorage.tabID)
      setIsConnected(true)
    }

    function onDisconnect() {
      console.log("DISCONNECTED")
      setIsConnected(false)
    }

    function onNewPacket(pkt) {
      setPacket(JSON.parse(pkt.data))
    }

    console.log("Setting Handlers")
    socket.onopen = onConnect
    socket.onclose = onDisconnect
    socket.onmessage = onNewPacket
  }, [])

  useEffect(() => {
    if (packet === undefined) {
      return
    }
    console.log("SETTINGS______________________________________________")
    console.log(settings)
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

  useEffect(() => {
    console.log("ROOM_____________________________________")
    console.log(settings.room)
    console.log("CLIENT__________________________________________")
    console.log(settings.client)
    console.log("GAME_____________________________________________")
    console.log(settings.game)
  }, [settings])

  return (
    <div className="App">
      <h1>test</h1>
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