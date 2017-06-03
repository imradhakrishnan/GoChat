package main

import (
	"log"
	"strings"
)

type Server struct {
	connRequests chan *Client     //for maintaining incoming requests
	exitRequests chan *Client     //for maintaining leaving requests
	clients      map[*Client]bool //map to maintain active clients
	messages     chan []byte      //to track messages from clients
}

func initServer() *Server {
	return &Server{
		make(chan *Client),
		make(chan *Client),
		make(map[*Client]bool),
		make(chan []byte),
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.connRequests:
			for members := range s.clients {
				client.send <- []byte("NAME:" + members.name)
			}
			s.clients[client] = true
		case client := <-s.exitRequests:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case msg := <-s.messages:
			msgarray := strings.Split(string(msg[:]), " ")
			log.Println(msgarray)
			from := strings.Split(msgarray[0], ":")
			var to []string
			if from[0] != "NAME" {
				to = strings.Split(msgarray[1], ":")
				log.Println(from, to)
			}
			for client := range s.clients {
				if from[0] == "NAME" || client.name == from[1] || client.name == to[1] {
					client.send <- msg
				}
			}
		}
	}
} //end of run
