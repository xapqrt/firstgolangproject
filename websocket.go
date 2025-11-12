//websocket handler
//upgrades http to websocket connection

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//upgrader config
//checkorigin returns true so cors doesnt mess things up

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


//handles websocket connections

func wsHandler(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("new websocket connection attempt from:", r.RemoteAddr)
	
	var conn *websocket.Conn
	var err error
	
	conn, err = upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		fmt.Println("upgrade failed:", err)
		return
	}
	
	client_counter = client_counter + 1
	var client_id string
	client_id = fmt.Sprintf("client-%d", client_counter)
	
	fmt.Println("client connected, id:", client_id)
	
	
	client := &Client{
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	
	client.hub.register <- client
	
	
	go client.writePump()
	go client.readPump()
	
	


	//lets js send it
	go sendInitialTasks(client)
	
	fmt.Println("client setup complete")
}







