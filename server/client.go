package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	tabID    string
	id       string
	roomCode string // "" if not in room
	name     string
	imguuid  string
	hub      *Hub
	conn     *websocket.Conn
	send     chan Packet
}

type ClientJSON struct {
	Id       string `json:"id"`
	RoomCode string `json:"roomCode"`
	Name     string `json:"name"`
	Imguuid  string `json:"imguuid"`
}

func (h *Hub) SendToClient(packet Packet, client *Client) {
	log.Println("Sending To Client")
	client.send <- packet
}

func (h *Hub) GetClientJSON(client *Client) string {
	cJSON := h.GetClientJSONStruct(client)
	dat, err := json.Marshal(cJSON)
	if err != nil {
		log.Println("Couldnt Parse Client JSON")
		return ""
	}
	return string(dat)
}

func (h *Hub) GetClientJSONStruct(client *Client) ClientJSON {
	cJSON := ClientJSON{
		Id:       client.id,
		RoomCode: client.roomCode,
		Name:     client.name,
		Imguuid:  client.imguuid,
	}
	return cJSON
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

func (c *Client) UpdateNonVitals(oC *Client) {
	c.roomCode = oC.roomCode
	c.name = oC.name
	c.imguuid = oC.imguuid
}

func (h *Hub) ClientDisconnect(client *Client) {
	h.clients[client] = false
	if client.roomCode != "" {
		client.roomCode = h.LeaveRoom(client)
	}
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

func Chk(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: Chk,
}

func (c *Client) ReadPackets() {
	defer func() {
		log.Println("Fuck Off Read my ASS packet")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, packetJson, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error ReadPackets(0) %v", err)
			}
			break
		}
		var packet Packet
		err = json.Unmarshal(packetJson, &packet)
		log.Println("Received Packet From Client")
		if err != nil {
			log.Printf("Error ReadPackets(1) %v", err)
		}
		if packet.Type == "setup" {
			sysCmd := strings.Split(packet.Data, " ")
			switch sysCmd[0] {
			case "clientconnect":
				if c.hub.ClientExists(sysCmd[1]) {
					c.hub.ClientRejoin(sysCmd[1], c)
				} else {
					log.Println("New Client " + sysCmd[1])
					c.tabID = sysCmd[1]
					c.hub.SendClientJSON(c)
				}
			}
			continue
		}
		c.hub.broadcast <- packet
	}
}

func (c *Client) WritePackets() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case packet, ok := <-c.send:
			if !ok {
				log.Println("Error WritePackets(0)")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(packet)
			if err != nil {
				log.Println("Error WritePackets(1)")
				log.Println(err)
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("New Client Joined - ClientID:%s\n", strconv.Itoa(nextClientID))
	client := &Client{
		id:       strconv.Itoa(nextClientID),
		roomCode: "",
		hub:      hub,
		conn:     conn,
		send:     make(chan Packet),
	}
	nextClientID += 1
	client.hub.register <- client

	go client.WritePackets()
	go client.ReadPackets()
}
