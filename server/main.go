package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
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

func uploadImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("File Uploading")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fx := strings.ReplaceAll(string(body), "=", "")
	dec, err := base64.RawStdEncoding.DecodeString(fx[strings.IndexByte(fx, ',')+1:])
	if err != nil {
		log.Printf("E0 %s", err.Error())
		return
	}
	uid := uuid.New().String()
	fileName := "images/" + uid + ".jpeg"
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("E1 %s", err.Error())
		return
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		log.Printf("E2 %s", err.Error())
		return
	}
	if err := f.Sync(); err != nil {
		log.Printf("E3 %s", err.Error())
		return
	}

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "%s", uid)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)

	f, err := os.ReadFile("./images/" + vars["uuid"] + ".jpeg")
	if err != nil {
		log.Println("Fialed Reading FIle")
		log.Println(err)
	}

	bs64 := base64.RawStdEncoding.EncodeToString(f)
	fmt.Fprintf(w, "data:image/jpeg;base64,%s", bs64)
}

var nextClientID int

func main() {
	fmt.Println("Drunk Games Server Running...")

	nextClientID = 1
	hub := NewHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/api/upload", uploadImage)
	r.HandleFunc("/api/image/get/{uuid}", getImage)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		log.Println("Someone Connected to ws")
		serveWs(hub, w, r)
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./build")))

	http.ListenAndServe(":80", r)
}
