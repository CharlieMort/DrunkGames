package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
)

type Packet struct {
	From string `json:"from"` //ClientID who sent msg - 0 if from server
	To   string `json:"to"`   //Recipent
	Type string `json:"type"` //Type of packet
	Data string `json:"data"` //The actual msg of the data
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string][]*Client
	broadcast  chan Packet
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Packet),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string][]*Client),
	}
}

func (h *Hub) CreateRoom() string {
	roomCode := GetRandomRoomCode()
	h.rooms[roomCode] = make([]*Client, 0)
	log.Printf("Created The Room:%s\n", roomCode)
	return roomCode
}

func (h *Hub) JoinRoom(client *Client, roomCode string) {
	type RoomJoinDataPacket struct {
		RoomCode string       `json:"roomCode"`
		Clients  []ClientJSON `json:"clients"`
	}
	if _, ok := h.rooms[roomCode]; ok {
		h.rooms[roomCode] = append(h.rooms[roomCode], client)
		client.roomCode = roomCode
		cJSON := make([]ClientJSON, 0)
		for _, client := range h.rooms[roomCode] {
			cJSON = append(cJSON, ClientJSON{
				Name:     client.name,
				Id:       client.id,
				RoomCode: client.roomCode,
			})
		}
		dat, err := json.Marshal(RoomJoinDataPacket{
			RoomCode: roomCode,
			Clients:  cJSON,
		})
		if err != nil {
			log.Println("Error Creating RoomJoinDataPacket")
		}

		h.SendClientJSON(client)
		h.SendPacket(Packet{
			From: "0",
			To:   roomCode,
			Type: "toRoom",
			Data: string(dat),
		})
		log.Printf("Client:%s Joined the Room:%s", client.id, roomCode)
	} else {
		h.SendPacket(Packet{
			From: "0",
			To:   client.id,
			Type: "error",
			Data: "The Room:" + roomCode + " Does Not Exist",
		})
		log.Printf("Client:%s Failed to join the Room:%s", client.id, roomCode)
	}
}

func (h *Hub) LeaveRoom(client *Client) {
	log.Println("Removing From Room " + client.roomCode)
	rIdx := slices.Index(h.rooms[client.roomCode], client)
	h.rooms[client.roomCode] = slices.Delete(
		h.rooms[client.roomCode],
		rIdx,
		rIdx+1,
	)
}

func (h *Hub) ShutDownClient(client *Client) {
	log.Printf("Client Disconnected ClientID:%s\n", client.id)
	if client.roomCode != "" {
		h.LeaveRoom(client)
	}
	delete(h.clients, client)
	close(client.send)
}

func (h *Hub) GetClientFromID(ID string) *Client {
	for client := range h.clients {
		if client.id == ID {
			return client
		}
	}
	return nil
}

func PrintClient(client Client) {
	fmt.Printf("ID:%s\nName:%s\nRoomCode:%s\n", client.id, client.name, client.roomCode)
}

func (h *Hub) SystemPacket(packet Packet) {
	log.Printf("Recieved Packet from Client From:%s\nData:%s\n", packet.From, packet.Data)
	sysCmd := strings.SplitN(packet.Data, " ", 2)
	switch sysCmd[0] {
	case "createroom":
		client := h.GetClientFromID(packet.From)
		if client == nil {
			log.Printf("Client:" + packet.From + " Couldnt Be Found")
			return
		}
		roomCode := h.CreateRoom()
		h.JoinRoom(client, roomCode)
	case "joinroom":
		client := h.GetClientFromID(packet.From)
		if client == nil {
			log.Printf("Client:" + packet.From + " Couldnt Be Found")
			return
		}
		roomCode := sysCmd[1]
		h.JoinRoom(client, roomCode)
	case "setclientname":
		client := h.GetClientFromID(packet.From)
		if client == nil {
			log.Printf("Client:" + packet.From + " Couldnt Be Found")
			return
		}
		client.name = sysCmd[1]
		h.SendClientJSON(client)
	}
}

func (h *Hub) SendToClient(packet Packet, client *Client) {
	log.Println("Sending To Client")
	select {
	case client.send <- packet:
	default:
		h.ShutDownClient(client)
	}
}

func (h *Hub) SendToRoom(packet Packet, roomCode string) {
	log.Println("Sending To Room")
	for _, client := range h.rooms[roomCode] {
		h.SendToClient(packet, client)
	}
}

func (h *Hub) SendToVacants(packet Packet) {
	log.Println("Sending To Vacants")
	for client := range h.clients {
		if client.roomCode == "" {
			h.SendToClient(packet, client)
		}
	}
}

func (h *Hub) SendToAll(packet Packet) {
	log.Println("Sending To All")
	for client := range h.clients {
		h.SendToClient(packet, client)
	}
}

func (h *Hub) SendPacket(packet Packet) {
	switch packet.Type {
	case "toClient", "error":
		client := h.GetClientFromID(packet.To)
		h.SendToClient(packet, client)
	case "toRoom":
		roomCode := packet.To
		if _, ok := h.rooms[roomCode]; ok {
			h.SendToRoom(packet, roomCode)
		} else {
			client := h.GetClientFromID(packet.From)
			h.SendToClient(Packet{
				From: "0",
				To:   client.id,
				Type: "error",
				Data: "Room: " + roomCode + " does not exist",
			}, client)
		}
	case "toVacants":
		h.SendToVacants(packet)
	case "toAll":
		h.SendToAll(packet)
	case "toSystem":
		h.SystemPacket(packet)
	default:
		log.Printf("Failed to send Packet\nFrom: %s\nTo: %s\nData: %s\n", packet.From, packet.To, packet.Data)
	}
}

func (h *Hub) GetClientJSON(client *Client) string {
	cJSON := ClientJSON{
		Id:       client.id,
		RoomCode: client.roomCode,
		Name:     client.name,
	}
	dat, err := json.Marshal(cJSON)
	if err != nil {
		log.Println("Couldnt Parse Client JSON")
		return ""
	}
	return string(dat)
}

func (h *Hub) SendClientJSON(client *Client) {
	dat := h.GetClientJSON(client)
	h.SendPacket(Packet{
		From: "0",
		To:   client.id,
		Type: "toClient",
		Data: dat,
	})
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.SendClientJSON(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.ShutDownClient(client)
			}
		case packet := <-h.broadcast:
			h.SendPacket(packet)
		}
	}
}
