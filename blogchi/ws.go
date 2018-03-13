package main

import (
	"log"
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
				removeConnection(u)
				log.Println("echoWS broadcast", err)
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
		m := "Unable to upgrade to websockets."
		log.Println(m, err)
		response.Send(w, 400, err.Error()).JSON()
		return
	}
	uuid := addConnection(conn)
	err = conn.WriteJSON(message{Type: "server:hello", Data: "Hello", UUID: uuid})
	if err != nil {
		removeConnection(uuid)
		log.Println("echoWS conn.WriteJSON", err)
		response.Send(w, 500, err.Error()).JSON()
		return
	}
	for {
		message := &message{}
		err := conn.ReadJSON(message)
		if err != nil {
			removeConnection(uuid)
			log.Println("echoWS for conn.ReadJSON", err)
			response.Send(w, 400, err.Error()).JSON()
			return
		}
		broadcast(message)

		// err = conn.WriteJSON(message)
		// if err != nil {
		// 	removeConnection(uuid)
		// 	log.Println("echoWS for conn.WriteJSON", err)
		// 	response.Send(w, 500, err.Error()).JSON()
		// 	return
		// }
		// messageType, p, err := conn.ReadMessage()
		// if err != nil {
		// 	if err == io.EOF {
		// 		log.Println("Websocket closed!")
		// 	} else {
		// 		log.Println("Error reading websocket message.", err)
		// 	}
		// 	return
		// }
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println("Error writing websocket message.", err)
		// 	return
		// }
	}
}
