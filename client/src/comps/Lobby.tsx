import React from "react"
import { IClient, IRoom } from "../types"
import RemoteImage from "./RemoteImage.tsx"

interface ILobbyProps {
    client: IClient
    room: IRoom
}

const Lobby = ({client, room}: ILobbyProps) => {
    return(
        <div>
            <h2>RoomCode: {room.roomCode}</h2>
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