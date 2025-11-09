//websocket client stuff
//each browser that connects gets one of these i think


package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

//client represents a single websocket connection
//not sure if i need all these fields but tutorial had them

type Client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}


var client_counter int = 0





//reads messages from websocket connection
//this runs forever until connection closes i guess

func (c *Client) readPump() {

	//cleanup when function exits
	defer func() {
		fmt.Println("client disconnecting...")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	//setting read deadline so it doesnt hang
	//60 seconds seems reasonable?
	var deadline time.Time
	deadline = time.Now().Add(60 * time.Second)
	c.conn.SetReadDeadline(deadline)

	//infinite loop to keep reading
	for {
		
		//read message from websocket
		messageType, message, err := c.conn.ReadMessage()
		
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		fmt.Println("got message:", string(message))
		fmt.Println("message type:", messageType)

		//send to hub so it can broadcast to everyone
		c.hub.broadcast <- message
	}
}
}




//writes messages to websocket connection
//sending stuff back to the browser

func (c *Client) writePump() {

	//close connection when done
	defer func() {
		c.conn.Close()
		fmt.Println("write pump closed")
	}()

	//keep listening for messages to send
	for {
		select {
		case message, ok := <-c.send:

			//check if channel is still open
			if !ok {
				//hub closed the channel so close connection
				fmt.Println("send channel closed")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			//write message to websocket
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