package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

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

	Connections = make([]Connection, 0)

	m sync.Mutex

	count = 0
)

type Connection struct {
	conn         *websocket.Conn
	Connected    bool
	writeChannel chan models.Message
	index        int
	profile      string
}

func newConnection(c *websocket.Conn) Connection {
	m.Lock()
	count++
	connection := Connection{
		conn:         c,
		Connected:    true,
		writeChannel: make(chan models.Message),
		index:        count,
	}
	Connections = append(Connections, connection)
	defer m.Unlock()
	return connection
}

func remove(s []Connection, index int) []Connection {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func (c *Connection) close() {
	m.Lock()
	index := -1
	for i, conn := range Connections {
		if conn.index == c.index {
			index = i
		}
	}
	if index > -1 {
		Connections = remove(Connections, index)
	}
	defer m.Unlock()
	c.Connected = false
	c.conn.Close()
}

// SendMessage ssending a message as a string
func SendMessage(message models.Message) {
	for _, conn := range Connections {
		if conn.Connected {
			if strings.EqualFold(message.Profile, conn.profile) {
				conn.writeChannel <- message
			}
		}
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
	conn := newConnection(c)
	defer conn.close()

	go readMessage(&conn)

	for {
		message := <-conn.writeChannel
		json, err := json.Marshal(message)
		if err != nil {
			clog.Logger.Errorf("json error: %v", err)
			continue
		}
		//clog.Logger.Infof("sending a message to the client: %s", message)
		err = c.WriteMessage(websocket.TextMessage, []byte(json))
		if err != nil {
			clog.Logger.Errorf("write error: %v", err)
			break
		}
	}
	conn.Connected = false
}

func readMessage(conn *Connection) {
	for {
		if conn.Connected {
			mt, message, err := conn.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				(*conn).close()
				break
			}
			var myMessage models.Message
			err = json.Unmarshal(message, &myMessage)
			if err != nil {
				clog.Logger.Errorf("json unmarshal: %v", err)
			}
			if myMessage.Command == "change" {
				if !strings.EqualFold(myMessage.Profile, conn.profile) {
					index := 0
					for i, lconn := range Connections {
						if lconn.index == conn.index {
							index = i
						}
					}
					Connections[index].profile = myMessage.Profile
				}
			}
			clog.Logger.Infof("recv: %d %s", mt, message)
		} else {
			break
		}
	}
}
