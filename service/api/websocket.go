package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	channel   = make(chan models.Message)
	connected = false
)

// SendMessage ssending a message as a string
func SendMessage(message models.Message) {
	if connected {
		channel <- message
	}
}

// ServeWs creates a new websocket connection to a requesting client
func ServeWs(w http.ResponseWriter, r *http.Request) {
	clog.Logger.Info("request websocket connection")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	connected = true
	defer c.Close()
	go readMessage(c)
	for {
		message := <-channel
		json, err := json.Marshal(message)
		if err != nil {
			clog.Logger.Errorf("json error: %v", err)
			continue
		}
		clog.Logger.Infof("sending a message to the client: %s", message)
		err = c.WriteMessage(websocket.TextMessage, []byte(json))
		if err != nil {
			clog.Logger.Errorf("write error: %v", err)
			break
		}
	}
	connected = false
}

func readMessage(c *websocket.Conn) {
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %d %s", mt, message)
	}
}
