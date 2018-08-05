package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
)

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool) //stores all the connected clients uniquely
var broadcast = make(chan Message) //message broadcast channel between all connected clients

func WebSocketHandler(ws *websocket.Conn) {
	clients[ws] = true //add new connection to clients
	for {
		var m Message
		websocket.JSON.Receive(ws, &m) //on receiving message from client
		broadcast <- m //broadcast it to channel
	}
}

func main() {
	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("public")) //set folder as file handler
	r.Handle("/", fs) //serve folder as file handler
	r.Handle("/ws", websocket.Handler(WebSocketHandler))
	go handleMessages() //run non blocking
	log.Println("Server started....")
	http.ListenAndServe(":8012", r)
}

//separate goroutine function to handle message broadcasting
func handleMessages() {
	for {//loop to send all messages found from broadcast channel to all clients
		msg := <- broadcast //get message from channel
		for client := range clients {
			websocket.JSON.Send(client, msg) //send to all clients
		}
	}
}