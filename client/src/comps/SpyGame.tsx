import React from "react"
import { ISpyGame } from "../types"

interface ISpyGameProps {
    gameData: ISpyGame
}

const SpyGame = ({gameData}: ISpyGameProps) => {
    return(
        <div>
            <h1>hehe spy game</h1>
        </div>
    )
}

export default SpyGame