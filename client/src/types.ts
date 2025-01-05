export interface IClient {
  id: string,
  roomCode: string,
  name: string,
  imguuid: string,
}

export interface IRoom {
  roomCode: string
  host: IClient
  clients: IClient[]
  gameType: string
}

export interface IPacket {
  from: string,
  to: string,
  type: string,
  data: string,
}

export interface ISpyGame {
  spies: number[]
  prompt: string
}

export interface ISettings {
  client: IClient | undefined
  room: IRoom | undefined
  game: ISpyGame | undefined
}