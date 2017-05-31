//complete reference from: https://github.com/gorilla/websocket/tree/master/examples/chat

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":12345", "http server port address")

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("home:", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	server := initServer()
	go server.run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleClient(server, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
