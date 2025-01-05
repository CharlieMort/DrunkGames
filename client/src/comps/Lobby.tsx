import React from "react"
import { IClient, IPacket, IRoom } from "../types"
import RemoteImage from "./RemoteImage.tsx"
import { socket } from "./socket.tsx"

interface ILobbyProps {
    client: IClient
    room: IRoom
}

const Lobby = ({client, room}: ILobbyProps) => {
    const StartGame = () => {
        socket.send(JSON.stringify({
            from: client.id,
            to: "0",
            type: "toSystem",
            data: `startgame spygame ${room.roomCode}`
        } as IPacket))
    }

    return(
        <div>
            <h2>RoomCode: {room.roomCode}</h2>
            {
                room.host.id === client.id && room.clients.length > 1
                ? <input className="bigButton" type="button" value="Start Game" onClick={StartGame} />
                : <></>
            }
            <div className="lobby">
                <div>
                    <h1>{client.name}</h1>
                    <RemoteImage uuid={client.imguuid} />
                </div>
                {
                    room.clients.map((cl) => {
                        if (cl.id != client.id) {
                            return (
                                <div>
                                    <h1>{cl.name}</h1>
                                    <RemoteImage uuid={cl.imguuid} />
                                </div>
                            )
                        }
                    })
                }
            </div>
        </div>
    )
}

export default Lobby