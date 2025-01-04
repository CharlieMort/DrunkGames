package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

const PORT = ":80"

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

	r := mux.NewRouter()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		log.Println("Someone Connected to ws")
		serveWs(hub, w, r)
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./build")))

	http.ListenAndServe(":80", r)
}
