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
	rooms      map[string]*Room
	broadcast  chan Packet
	register   chan *Client
	unregister chan *Client
}

type Room struct {
	RoomCode string    `json:"roomCode"`
	Host     *Client   `json:"host"`
	Game     Game      `json:"game"`
	Clients  []*Client `json:"clients"`
}

type RoomJSON struct {
	RoomCode string       `json:"roomCode"`
	Host     ClientJSON   `json:"host"`
	Clients  []ClientJSON `json:"clients"`
	GameType string       `json:"gameType"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Packet),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) ClientExists(tabID string) bool {
	for cl := range h.clients {
		if cl.tabID == tabID && cl.tabID != "" {
			return true
		}
	}
	return false
}

func (h *Hub) ClientRejoin(tabID string, currClient *Client) {
	for cl := range h.clients {
		if cl.tabID == tabID {
			log.Printf("Client Rejoined - MV: %s -> %s\n", currClient.id, cl.id)
			//h.ShutDownClient(currClient)
			currClient.UpdateNonVitals(cl)
			h.SendClientJSON(currClient)
			log.Println(currClient.roomCode)
			if currClient.roomCode != "" {
				h.JoinRoom(currClient, currClient.roomCode)
				h.LeaveRoom(cl)
			}
		}
	}
}

func (h *Hub) CreateRoom() string {
	roomCode := GetRandomRoomCode()
	h.rooms[roomCode] = &Room{
		RoomCode: roomCode,
		Host:     nil,
		Game:     nil,
		Clients:  make([]*Client, 0),
	}
	log.Printf("Created The Room:%s\n", roomCode)
	return roomCode
}

func (h *Hub) JoinRoom(client *Client, roomCode string) {
	roomCode = strings.ToLower(roomCode)
	if _, ok := h.rooms[roomCode]; ok {
		if h.rooms[roomCode].Host == nil {
			h.rooms[roomCode].Host = client
		}
		h.rooms[roomCode].Clients = append(h.rooms[roomCode].Clients, client)
		client.roomCode = roomCode
		h.SendClientJSON(client)
		h.SendRoomUpdate(roomCode)
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

func (h *Hub) LeaveRoom(client *Client) string {
	log.Println("Removing From Room " + client.roomCode)
	rIdx := slices.Index(h.rooms[client.roomCode].Clients, client)
	if rIdx == -1 {
		return ""
	}
	rc := client.roomCode
	h.rooms[client.roomCode].Clients = slices.Delete(
		h.rooms[client.roomCode].Clients,
		rIdx,
		rIdx+1,
	)
	if h.rooms[client.roomCode].Host == client {
		h.rooms[client.roomCode].Host = nil
		if len(h.rooms[client.roomCode].Clients) == 0 {
			delete(h.rooms, client.roomCode)
			client.roomCode = ""
			return rc
		} else {
			h.rooms[client.roomCode].Host = h.rooms[client.roomCode].Clients[0]
		}
	}
	h.SendRoomUpdate(client.roomCode)
	client.roomCode = ""
	return rc
}

func (h *Hub) SendRoomUpdate(roomCode string) {
	cJSON := make([]ClientJSON, 0)
	for _, client := range h.rooms[roomCode].Clients {
		cJSON = append(cJSON, h.GetClientJSONStruct(client))
	}
	gType := ""
	if h.rooms[roomCode].Game != nil {
		gType = h.rooms[roomCode].Game.GetType()
	}
	fmt.Println(gType)
	dat, err := json.Marshal(RoomJSON{
		RoomCode: roomCode,
		Host:     h.GetClientJSONStruct(h.rooms[roomCode].Host),
		Clients:  cJSON,
		GameType: gType,
	})
	if err != nil {
		log.Println("Error Creating RoomJoinDataPacket")
	}

	h.SendPacket(Packet{
		From: "0",
		To:   roomCode,
		Type: "toRoom",
		Data: string(dat),
	})
}

func (h *Hub) SystemPacket(packet Packet) {
	log.Printf("Recieved Packet from Client\nFrom:%s\nData:%s\n", packet.From, packet.Data)
	sysCmd := strings.Split(packet.Data, " ")
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
		client.name = strings.SplitN(packet.Data, " ", 2)[1]
		h.SendClientJSON(client)
	case "setclientimage":
		client := h.GetClientFromID(packet.From)
		if client == nil {
			log.Printf("Client:" + packet.From + " Couldnt Be Found")
			return
		}
		client.imguuid = sysCmd[1]
		h.SendClientJSON(client)
	case "startgame":
		roomCode := sysCmd[2]
		switch sysCmd[1] {
		case "spygame":
			log.Println("Startup SpyGame")
			h.rooms[roomCode].Game = &SpyGame{
				Stage:         0,
				Hub:           h,
				Room:          h.rooms[roomCode],
				Spies:         make([]*Client, 0),
				Prompt:        "",
				ReadyCount:    0,
				QuestionOrder: make([]*Client, 0),
				Votes:         make(map[int]int),
				Dead:          make([]*Client, 0),
			}
		}
		h.rooms[roomCode].Game.StartGame()
	}
}

func (h *Hub) SendToRoom(packet Packet, roomCode string) {
	log.Println("Sending To Room")
	for _, client := range h.rooms[roomCode].Clients {
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
	case "heartbeat":
		log.Printf("HES ALIVE %s", packet.To)
		client := h.GetClientFromID(packet.To)
		if client != nil {
			h.SendToClient(packet, client)
		}
	case "clientUpdate":
		client := h.GetClientFromID(packet.To)
		if client != nil {
			h.SendClientJSON(client)
		}
	case "toClient", "error":
		client := h.GetClientFromID(packet.To)
		if client != nil {
			h.SendToClient(packet, client)
		}
	case "toRoom", "toGame":
		roomCode := packet.To
		if _, ok := h.rooms[roomCode]; ok {
			h.SendToRoom(packet, roomCode)
		} else {
			client := h.GetClientFromID(packet.From)
			if client == nil {
				break
			}
			log.Println("Room Does Not Exist")
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

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("The Cunt disconnected")
			if _, ok := h.clients[client]; ok {
				h.ClientDisconnect(client)
			}
		case packet := <-h.broadcast:
			h.SendPacket(packet)
		}
	}
}
