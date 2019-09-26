package common

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var DefaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandleFunc func(clientID int, message types.H, args ...interface{})

type MessageHandlers interface {
	OnOpen(clientID int)
	OnClose(clientID int)
}

type WebSocket struct {
	upgrader           websocket.Upgrader
	clients            []*websocket.Conn
	messageHandlers    *MessageHandlers
	objMessageHandlers reflect.Value
	mutex              *sync.Mutex
}

func (ws *WebSocket) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	client, err := ws.upgrader.Upgrade(w, r, nil)
	ws.mutex.Lock()
	clientID := ws.addClient(client)
	ws.mutex.Unlock()
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	(*ws.messageHandlers).OnOpen(clientID)
	for {
		mt, message, err := client.ReadMessage()
		if err != nil {
			logger.HandleError(err)
			break
		}

		data := struct {
			Event string
			Data  types.H
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
			decorators := ws.objMessageHandlers.Elem().FieldByName("Decorators")

			if decorators.Kind() == reflect.Slice {
				for i := decorators.Len() - 1; i >= 0; i-- {
					if decorators.Index(i).Kind() == reflect.Func {
						method = decorators.Index(i).Call([]reflect.Value{method})[0]
					}
				}
			}

			method.Call([]reflect.Value{
				reflect.ValueOf(clientID),
				reflect.ValueOf(data.Data),
			})

		} else {
			body, err = json.Marshal(types.H{
				"event": "event_error",
				"data": types.H{
					"error":   true,
					"message": fmt.Sprintf("Event %s not supported", data.Event),
					"code":    404,
				},
			})
			logger.Logger.Warning(string(body))
			err = client.WriteMessage(mt, body)
			if err != nil {
				logger.HandleError(err)
				break
			}
			break
		}
	}

	defer func() {
		(*ws.messageHandlers).OnClose(clientID)
		ws.mutex.Lock()
		ws.removeClient(clientID)
		ws.mutex.Unlock()
		logger.HandleError(client.Close())
	}()
}

func (ws WebSocket) Emit(clientID int, event string, data types.H) {
	body, err := json.Marshal(types.H{
		"event": event,
		"data":  data,
	})
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	client := ws.Client(clientID)
	if client != nil {
		logger.HandleError(ws.Client(clientID).WriteMessage(1, body))
	}
}

func (ws *WebSocket) addClient(client *websocket.Conn) int {
	clientID := -1
	flagIsSetLast := false
	for i := len(ws.clients) - 1; i >= 0; i-- {
		if ws.clients[i] != nil {
			flagIsSetLast = true
			clientID = i - 1
		} else if !flagIsSetLast {
			ws.clients = append(ws.clients[:i], ws.clients[i+1:]...)
		} else {
			break
		}
	}
	if clientID == -1 {
		ws.clients = append(ws.clients, client)
		clientID = len(ws.clients) - 1
	} else {
		ws.clients[clientID] = client
	}
	return clientID
}

func (ws *WebSocket) removeClient(clientID int) {
	ws.clients[clientID] = nil
}

func (ws *WebSocket) Clients() []*websocket.Conn {
	return ws.clients
}

func (ws *WebSocket) Client(clientID int) *websocket.Conn {
	if len(ws.clients) > clientID {
		return ws.clients[clientID]
	}
	return nil
}

func (ws *WebSocket) Init(h MessageHandlers, upgrader websocket.Upgrader) {
	ws.clients = make([]*websocket.Conn, 0)
	ws.upgrader = upgrader
	ws.messageHandlers = &h
	ws.mutex = new(sync.Mutex)
	ws.objMessageHandlers = reflect.ValueOf(h)
}
