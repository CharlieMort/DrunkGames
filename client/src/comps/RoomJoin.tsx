import React, {useEffect, useState} from "react"
import { socket } from "./socket.tsx"
import { IClient, IPacket } from "../types.ts"
import Cam from "./Camera.jsx"

interface IRoomJoinProps {
    client: IClient
}

const RoomJoin = ({client}: IRoomJoinProps) => {
    const [name, setName] = useState<string>("")
    const [nameTmp, setNameTmp] = useState<string>("")
    const [roomCode, setRoomCode] = useState<string>("")
    const [joinRoom, setJoinRoom] = useState<boolean>(false)
    const [imgUUID, setImgUUID] = useState<string>("")

    const JoinRoom = () => {
        socket.send(JSON.stringify({
            from: client.id,
            to: "0",
            type: "toSystem",
            data: `joinroom ${roomCode}` 
        } as IPacket))
    }

    const SubmitName = () => {
        setName(nameTmp)
        socket.send(JSON.stringify({
            from: client.id,
            to: "0",
            type: "toSystem",
            data: `setclientname ${nameTmp}` 
        } as IPacket))
    }

    const CreateRoom = () => {
        socket.send(JSON.stringify({
            from: client.id,
            to: "",
            type: "toSystem",
            data: "createroom"
        }  as IPacket))
    }

    useEffect(() => {
        if (imgUUID !== "") {
            socket.send(JSON.stringify({
                from: client.id,
                to: "",
                type: "toSystem",
                data: `setclientimage ${imgUUID}`
            } as IPacket))
        }
    }, [imgUUID])

    return (
        <div className="panel">
            {
                name === "" && client.name === ""
                ? <div>
                    <h1>{sessionStorage.tabID}</h1>
                    <div className="panel-bot">
                        <label htmlFor="username">Enter Your Name</label>
                        <input className="bigTextInput" id="username" type="text" placeholder="Your Name" value={nameTmp} onChange={(e) => setNameTmp(e.target.value)}/>
                        <input className="bigButton" id="joinbutton" type="button" value="Next" onClick={SubmitName} />
                    </div>
                </div>
                : client.imguuid === ""
                    ? <div>
                        <Cam setImgUUID={setImgUUID} />
                    </div>
                    : joinRoom
                        ? <div className="panel-bot">
                            <label htmlFor="roomcode">Enter the Room Code</label>
                            <input className="bigTextInput" id="roomcode" type="text" placeholder="Room Code" value={roomCode} onChange={(e) => setRoomCode(e.target.value)} />
                            <input className="bigButton" id="joinbutton" type="button" value="Join" onClick={JoinRoom} />
                        </div>
                        : <div className="panel-bot">
                            <h1>{client.imguuid}</h1>
                            <input className="bigButton" id="joinRoomButton" type="button" value="Join Room" onClick={() => setJoinRoom(true)} />
                            <input className="bigButton" id="createRoomButton" type="button" value="Create Room" onClick={CreateRoom} />
                        </div>
            }
        </div>
    )
}

export default RoomJoin

{/* <div className="panel-top">
    <h2>Join a Room !</h2>
</div>
<div className="panel-bot">
    <label htmlFor="username">Enter Your Name</label>
    <input id="username" type="text" placeholder="Your Name" value={roomCode} onChange={(e) => setRoomCode(e.target.value)}/>
    <label htmlFor="roomcode">Enter the Room Code</label>
    <input id="roomcode" type="text" placeholder="Room Code" />
    <input id="joinbutton" type="button" value="Join" onClick={JoinRoom} />
</div> */}