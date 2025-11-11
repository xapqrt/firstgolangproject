//hub manages all websocket clients
//idk if this is the best way but tutorials do it like this

package main

import (
	"fmt"
)

//hub maintains active clients and broadcasts messages

type Hub struct {
	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

//creates new hub instance
//probably should only have one of these

func newHub() *Hub {

	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

//runs the hub in background
//infinite loop with select

func (h *Hub) run() {

	fmt.Println("hub is running...")

	for {
		select {

		case client := <-h.register:
			h.clients[client] = true
			var total int
			total = len(h.clients)
			fmt.Println("new client registered, total:", total)

		case client := <-h.unregister:
			
			_, exists := h.clients[client]
			
			if exists {
				delete(h.clients, client)
				close(client.send)
				fmt.Println("client unregistered, remaining:", len(h.clients))
			}

		case message := <-h.broadcast:

			client_count := len(h.clients)
			fmt.Println("broadcasting to", client_count, "clients")

			//send to all connected clients
			for client := range h.clients {

				select {
				case client.send <- message:
					//sent

				default:
					//channel full maybe
					close(client.send)
					delete(h.clients, client)
					fmt.Println("removed slow client")
				}
			}
		}
	}
}
