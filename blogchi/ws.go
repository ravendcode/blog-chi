package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var connections = map[string]*websocket.Conn{}

func addConnection(conn *websocket.Conn) string {
	uuid := getUUID()
	connections[uuid] = conn
	return uuid
}

func removeConnection(uuid string) {
	delete(connections, uuid)
}

func broadcast(m *message) {
	for uuid, conn := range connections {
		go func(m *message, u string, conn *websocket.Conn) {
			err := conn.WriteJSON(m)
			if err != nil {
				logger.Err(err)
				removeConnection(u)
				return
			}
		}(m, uuid, conn)
	}
}

func getUUID() string {
	u := uuid.NewV4()
	return u.String()
}

type message struct {
	UUID string      `json:"uuid"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func echoWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err(err)
		response.Send(w, 400, err).JSON()
		return
	}
	uuid := addConnection(conn)
	err = conn.WriteJSON(message{Type: "server:hello", Data: "Hello", UUID: uuid})
	if err != nil {
		logger.Err(err)
		removeConnection(uuid)
		return
	}
	for {
		message := &message{}
		err := conn.ReadJSON(message)
		if err != nil {
			if c, ok := err.(*websocket.CloseError); ok {
				if c.Code == 1001 {
					logger.Info("uuid:", uuid, err)
					removeConnection(uuid)
					break
				}
			}
			logger.Err(err)
			removeConnection(uuid)
			break
		}
		broadcast(message)
	}
}
