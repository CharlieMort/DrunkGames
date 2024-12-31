import React, {useState} from "react"
import { socket } from "./socket.tsx"
import { IClient } from "../types.ts"

interface IRoomJoinProps {
    client: IClient
}

const RoomJoin = ({client}: IRoomJoinProps) => {
    const [name, setName] = useState<string>("")
    const [nameTmp, setNameTmp] = useState<string>("")
    const [roomCode, setRoomCode] = useState<string>("")
    const [joinRoom, setJoinRoom] = useState<boolean>(false)

    const JoinRoom = () => {
        socket.send(JSON.stringify({
            "from": client.id,
            "to": "0",
            "type": "toSystem",
            "data": "joinroom "+roomCode
        }))
    }

    const SubmitName = () => {
        setName(nameTmp)
        socket.send(JSON.stringify({
            "from": client.id,
            "to": "0",
            "type": "toSystem",
            "data": "setclientname "+nameTmp  
        }))
    }

    const CreateRoom = () => {
        socket.send(JSON.stringify({
            "from": client.id,
            "to": "",
            "type": "toSystem",
            "data": "createroom"
        }))
    }

    return (
        <div className="panel">
            {
                name === ""
                ? <div>
                    <div className="panel-bot">
                        <label htmlFor="username">Enter Your Name</label>
                        <input className="bigTextInput" id="username" type="text" placeholder="Your Name" value={nameTmp} onChange={(e) => setNameTmp(e.target.value)}/>
                        <input className="bigButton" id="joinbutton" type="button" value="Next" onClick={SubmitName} />
                    </div>
                </div>
                : joinRoom
                    ? <div className="panel-bot">
                        <label htmlFor="roomcode">Enter the Room Code</label>
                        <input className="bigTextInput" id="roomcode" type="text" placeholder="Room Code" value={roomCode} onChange={(e) => setRoomCode(e.target.value)} />
                        <input className="bigButton" id="joinbutton" type="button" value="Join" onClick={JoinRoom} />
                    </div>
                    : <div className="panel-bot">
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