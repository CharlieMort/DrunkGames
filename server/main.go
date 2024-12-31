package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func GetRandomRoomCode() string {
	// REMOVED CHARS - c i 1 l j g
	rnd := "abdefhkmnopqrstuvqxyz023456789"
	lng := 5
	out := ""
	for i := 0; i < lng; i++ {
		out = out + string(rnd[rand.Intn(len(rnd))])
	}
	return out
}

var nextClientID int

func main() {
	fmt.Println("Drunk Games Server Running...")

	nextClientID = 1
	hub := NewHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		log.Println("Someone Connected to /")
		w.Write([]byte("Thanks for the req"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		log.Println("Someone Connected to ws")
		serveWs(hub, w, r)
	})

	fmt.Println("Listening On Port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
