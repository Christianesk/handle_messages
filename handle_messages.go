package handle_messages

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Greeting string `json:"greeting"`
}

var (
	wsUpgrader = websocket.Upgrader{}

	wsConn *websocket.Conn
)

func SendMessageToGameClient(msg string) {
	fmt.Printf("[CURRENT MESSAGE TO CLIENT]: %s", msg)
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Print(err)
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		// check the http.Request
		// make sure it's OK to access
		return true
	}
	var err error
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	// event loop
	for {
		var msg Message

		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error reading JSON: %s\n", err.Error())
			break
		}

		fmt.Printf("Message Received: %s\n", msg.Greeting)
		SendMessage(msg.Greeting)
	}
}

func SendMessage(msg string) {
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
	}
}

func StartWebSocketToClient() {
	router := mux.NewRouter()

	router.HandleFunc("/socket", WsEndpoint)

	log.Fatal(http.ListenAndServe(":9100", router))
}
