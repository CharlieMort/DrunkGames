import React from "react"
import { IClient, IRoom } from "../types"

interface ILobbyProps {
    client: IClient
    room: IRoom
}

const Lobby = ({client, room}: ILobbyProps) => {
    return(
        <div>
            <h2>RoomCode: {room.roomCode}</h2>
            <h1>{client.name}</h1>
            {
                room.clients.map((cl) => {
                    if (cl.id != client.id) {
                        return (
                            <h1>{cl.name}</h1>
                        )
                    }
                })
            }
        </div>
    )
}

export default Lobby