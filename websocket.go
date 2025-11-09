//websocket handler//alr lets get to the websocket part for the realtime updates and shit

//upgrades http to websocket connection



package main



import (package main

	"fmt"

	"net/http"import (


	"github.com/gorilla/websocket"

)	"fmt"

	"net/http"

	"sync"

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,

	WriteBufferSize: 1024,	"github.com/gorilla/websocket"

	)

	CheckOrigin: func(r *http.Request) bool {

		return true

	},

}

//upgrader config (found this in broscode tutorial)



//handles websocket connections



func wsHandler(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{

	fmt.Println("new websocket connection attempt")

	

	conn, err := upgrader.Upgrade(w, r, nil)

		ReadBufferSize: 1024,

	if err != nil {	WriteBufferSize: 1024,

		fmt.Println("upgrade failed:", err)

		return

	}

	

	client_counter = client_counter + 1

	fmt.Println("client connected, id:", client_counter)

	

		CheckOrigin: func(r *http.Request) bool {

	client := &Client{

		hub: hub,		return true

		conn: conn,	},

		send: make(chan []byte, 256),}

	}

	

	client.hub.register <- client

	

	

	go client.writePump()

	//representing a single user?

	go client.readPump()

}



type Client struct {

	conn *websocket.conn
	id string
	send chan []byte    

}





type Hub struct {



	clients map[*Client]bool

	broadcast chan[]byte




	register chan *Client
	unregister chan *Client

	mu sync.Mutex

}




//tryna make a global hub



var hub *Hub



//func

func newHub() *Hub {

	return &Hub {

		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}






//backgrond bug hub


func (h *Hub) run(){

	fmt.Println("hub starting to listen for events...")



	for{

		select {


		case client := <-h.register:

			h.mu.Lock()
			h.clients[client] = true


				h.mu.Unlock()


				fmt.Println("client connected: ", client.id)
				fmt.Println("total clients: ", len(h.clients))

		case client := <-h.unregister:



			h.mu.Lock()



			if _, ok := h.clients[client]; ok {

				delete(h.clients, client)
				close(client.send)


				fmt.Println("client disconnected: ", client.id)
				fmt.Println("total clients: ", len(h.clients))

			}

			h.mu.Unlock()


		case message := <-h.broadcast:



			h.mu.Lock()


			//send msg to connected users if dead





			for client := range h.clients {

				select {

				case client.send <- message:


				default: 

									close(client.send)
									delete(h.clients, client)
									fmt.Println("removed stuck client:", client.id)
				}
			}



			//lets hope it works


			h.mu.Unlock()
		}
	}
}


//reads messages from websocket connection


func (c *Client) readPump() {


	defer func() {

		hub. unregister <- c
		c.conn.Close()
	}()



	for {

		_, message, err := c.conn.ReadMessage()


		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {




				fmt.Println("read error: ", err)

			}
					break

		}

		fmt.Println("received message from ", c.id, ":", string(message))


		//broadcast msg

		hub.broadcast <- message

	}
}


//writing msg to websocket connection

//bro nah bro this shit is toughh
func (c * Client) writePump() {

	defer func() {

		c.conn.Close()
	}()




	for {


		message, ok := <-c.send
		if !ok {


			c.conn.WriteMessage(websocket. CloseMessage, []byte{})
			return
		}
	}
}







//adding a breakpoint


//i think this is the part where itll break




//handling webscoket upgrade and client setup

func wsHandler(w http.ResponseWriter, r *http.Request) {


	fmt.Println("websocket connection attempt from: ", r.RemoteAddr)



	//upgrading http to websocket



	conn, err := upgrader.Upgrade(w, r, nil)



	if err !=nil {

		fmt.Println("Upgrade FAiled: ", err)
		return

	}


	//create new client


	var client_id string
	client_id = fmt.Sprintf("client-%d", len(hub.clients)+1)



	client := &Client{
		conn: conn,
		id: client_id,
		send: make(chan []byte, 256),


	}


	hub.register <- client


	//background tasks as they ri imp

	go client.writePump()
	go client.readPump()


	fmt.Println("client setup complete: ", client_id)


}






