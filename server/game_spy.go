package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
)

type SpyGame struct {
	Stage         int
	Room          *Room
	Hub           *Hub
	Prompt        string
	NumOfSpies    int
	Spies         []*Client
	ReadyCount    int
	QuestionOrder []*Client
	Votes         map[int]int
	Dead          []*Client
}

type SpyGameJSON struct {
	Spies  []int  `json:"spies"`
	Prompt string `json:"prompt"`
}

func GetRandomPrompt() string {
	prompts := []string{
		"Beach",
		"Hotel",
		"University",
		"Pub",
		"Weatherspoons",
		"America",
		"India",
		"Rave",
		"Mickey Mouse Clubhouse",
	}
	return prompts[rand.IntN(len(prompts))]
}

func (g *SpyGame) StartGame() {
	log.Println("Startubg SoyGane")
	g.Hub.SendRoomUpdate(g.Room.RoomCode)
	g.Stage = 1
	g.ReadyCount = 0
	sIdx := make([]int, 0)
	for i := 0; i < g.NumOfSpies; i++ {
		idx := rand.IntN(len(g.Room.Clients))
		g.Spies = append(g.Spies, g.Room.Clients[idx])
		sIdx = append(sIdx, idx)
	}
	g.Prompt = GetRandomPrompt()

	cJSON := SpyGameJSON{
		Spies:  sIdx,
		Prompt: g.Prompt,
	}

	dat, err := json.Marshal(cJSON)

	if err != nil {
		log.Println("Fuck me spy game didnt work start game whatever" + err.Error())
	}

	g.Hub.SendPacket(Packet{
		From: "0",
		To:   g.Room.RoomCode,
		Type: "toGame",
		Data: string(dat),
	})
}

func (g *SpyGame) GetType() string {
	log.Println("Retyrbubg SPYGAME AS TYPE")
	return "spygame"
}
