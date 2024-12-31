package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	roomCode string // "" if not in room
	name     string
	hub      *Hub
	conn     *websocket.Conn
	send     chan Packet
}

type ClientJSON struct {
	Id       string `json:"id"`
	RoomCode string `json:"roomCode"`
	Name     string `json:"name"`
}

func Chk(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: Chk,
}

func (c *Client) ReadPackets() {
	defer func() {
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
		if err != nil {
			log.Printf("Error ReadPackets(1) %v", err)
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
				log.Fatal(err)
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
	client := &Client{id: strconv.Itoa(nextClientID), roomCode: "", hub: hub, conn: conn, send: make(chan Packet)}
	nextClientID += 1
	client.hub.register <- client

	go client.WritePackets()
	go client.ReadPackets()
}
