export interface IClient {
  id: string,
  roomCode: string,
  name: string
}

export interface IRoom {
  roomCode: string
  clients: IClient[]
}

export interface IPacket {
  from: string,
  to: string,
  type: string,
  data: string,
}