//websocket client stuff
//each browser that connects gets one of these i think


package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)


type Client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}


var client_counter int = 0




func (c *Client) readPump() {

	defer func() {
		fmt.Println("client disconnecting...")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	//setting read deadline so it doesnt hang
	var deadline time.Time
	deadline = time.Now().Add(60 * time.Second)
	c.conn.SetReadDeadline(deadline)

	for {
		
		messageType, message, err := c.conn.ReadMessage()
		
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		fmt.Println("got message:", string(message))
			fmt.Println("message type:", messageType)

		handleClientMessage(message)
	}
}








func (c *Client) writePump() {

	defer func() {
		c.conn.Close()
		fmt.Println("write pump closed")
	}()

	for {
		select {
		case message, ok := <-c.send:

			if !ok {
				fmt.Println("send channel closed")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			var err error
			err = c.conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				fmt.Println("write error:", err)
				return
			}
			
				fmt.Println("sent message to client")

		}
	}
}