import React from "react"
import { IRoom, ISpyGame } from "../types"
import SpyGame from "./SpyGame.tsx"

interface IGameRouterProps {
    room: IRoom
    gameData: ISpyGame
}

const GameRouter = ({room, gameData}: IGameRouterProps) => {
    switch(room.gameType) {
        case "spygame":
            return <SpyGame gameData={gameData} />
    }

    return <></>
}

export default GameRouter