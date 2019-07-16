package common

import (
	"encoding/json"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.HandleError(err)
			break
		}

		data := struct {
			Event string                 `json:"event"`
			Data  map[string]interface{} `json:"data"`
		}{}

		// TODO switch case Event

		err = json.Unmarshal(message, &data)
		if err != nil {
			logger.HandleError(err)
			break
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.HandleError(err)
			break
		}
	}

	defer logger.HandleError(c.Close())
}
