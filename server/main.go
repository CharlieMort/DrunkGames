package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const PORT = ":80"
const DEBUG = true

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
	var fileName string
	if DEBUG {
		fileName = "./images/" + uid + ".jpeg"
	} else {
		fileName = "/go/bin/images/" + uid + ".jpeg"
	}
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

	var fileName string
	if DEBUG {
		fileName = "./images/" + vars["uuid"] + ".jpeg"
	} else {
		fileName = "/go/bin/images/" + vars["uuid"] + ".jpeg"
	}
	f, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Fialed Reading FIle")
		log.Println(err)
	}

	bs64 := base64.RawStdEncoding.EncodeToString(f)
	fmt.Fprintf(w, "data:image/jpeg;base64,%s", bs64)
}

var nextClientID int

type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	path := filepath.Join(h.staticPath, r.URL.Path)
	indexPath := filepath.Join(h.staticPath, h.indexPath)
	log.Printf("URL: %s outPATH: %s indexPath:%s\n", r.URL.Path, path, indexPath)
	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		log.Println("FILE DONT EXIST")
		http.ServeFile(w, r, indexPath)
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

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
	spa := spaHandler{staticPath: "/go/bin/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.ListenAndServe(PORT, r)
}
