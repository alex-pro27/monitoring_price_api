package common

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var DefaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageHandlers interface {
	OnOpen(clientID int)
	OnClose(clientID int)
}

type WebSocket struct {
	upgrader websocket.Upgrader
	clients []*websocket.Conn
	mt map[int]int
	messageHandlers MessageHandlers
	objMessageHandlers reflect.Value
}

func (ws *WebSocket) Handle(w http.ResponseWriter, r *http.Request)  {
	var err error
	client, err := ws.upgrader.Upgrade(w, r, nil)
	clientID := ws.addClient(client)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	for {
		mt, message, err := client.ReadMessage()
		ws.mt[clientID] = mt
		if err != nil {
			logger.HandleError(err)
			break
		}

		data := struct {
			Event string                 `json:"event"`
			Data  map[string]interface{} `json:"data"`
		}{}

		err = json.Unmarshal(message, &data)
		if err != nil {
			logger.HandleError(err)
			break
		}

		var body []byte
		pattern := regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
		event := "On" + pattern.ReplaceAllStringFunc(data.Event, func(s string) string {
			return strings.ToUpper(strings.Replace(s, "_", "", -1))
		})

		method := ws.objMessageHandlers.MethodByName(event)
		if method.Kind() != reflect.Invalid {
			method.Call([]reflect.Value{
				reflect.ValueOf(clientID),
				reflect.ValueOf(data),
			})[0].Interface()
		} else {
			body = []byte(fmt.Sprintf("Event %s not supported", data.Event))
			logger.Logger.Warning(body)
			err = client.WriteMessage(mt, body)
			if err != nil {
				logger.HandleError(err)
				break
			}
			break
		}
	}

	defer func() {
		logger.HandleError(client.Close())
		ws.messageHandlers.OnClose(clientID)
		ws.removeClient(clientID)
	}()
}

func (ws WebSocket) Emmit(clientID int, data map[string]interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	logger.HandleError(ws.clients[clientID].WriteMessage(ws.mt[clientID], body))
}

func (ws *WebSocket) SetUpgrader(upgrader websocket.Upgrader)  {
	ws.upgrader = upgrader
}

func (ws *WebSocket) addClient(client *websocket.Conn) int {
	ws.clients = append(ws.clients, client)
	return len(ws.clients) - 1
}

func (ws *WebSocket) removeClient(clientID int) {
	ws.clients = append(ws.clients[:clientID], ws.clients[clientID + 1:]...)
}

func (ws *WebSocket) Clients() []*websocket.Conn  {
	return ws.clients
}

func (ws *WebSocket) Client(clientID int) *websocket.Conn  {
	return ws.clients[clientID]
}

func (ws *WebSocket) SetEventHandlers(h MessageHandlers)  {
	ws.messageHandlers = h
	ws.objMessageHandlers = reflect.ValueOf(h)
}

