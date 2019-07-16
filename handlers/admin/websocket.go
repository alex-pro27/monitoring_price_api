package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

type AdminWebSocketHandler struct {
	Server *common.WebSocket
	Users map[string]int
	db *gorm.DB
}

func (h AdminWebSocketHandler) Emit(token, event string, message types.H)  {
	clientID := AdminWebSocket.Users[token]
	h.Server.Emit(clientID, event, message)
}

func (h *AdminWebSocketHandler)OnOpen(clientID int) {
	h.Server.Emit(clientID, "on_open", types.H{
		"message": fmt.Sprintf("client on connected %d", clientID),
	})
}

func (h *AdminWebSocketHandler) OnClose(clientID int) {
	var token string
	for _token, id := range h.Users {
		if id == clientID {
			token = _token
			break
		}
	}
	if token != "" {
		delete(h.Users, token)
	}
}

func (h *AdminWebSocketHandler) OnConnect(clientID int, message types.H) {
	user := models.User{}
	token := message["token"]
	if token != nil {
		user.Manager(h.db).GetUserByToken(message["token"].(string))
	}
	if !user.IsStaff {
		h.Server.Emit(clientID, "on_connect", types.H{
			"error": true,
			"code": 403,
			"message": "Permission denied",
		})
		logger.Logger.Warning("Admin Websocket error: Permission denied")
		logger.HandleError(h.Server.Client(clientID).Close())
		return
	}
	h.Users[user.Token.Key] = clientID
	h.Server.Emit(clientID, "on_connect", types.H{
		"message": "Connected!",
	})
}

var AdminWebSocket *AdminWebSocketHandler

func StartWebsocket(w http.ResponseWriter, r *http.Request) {
	ws := common.WebSocket{}
	db := context.Get(r, "DB").(*gorm.DB)
	AdminWebSocket = &AdminWebSocketHandler{
		Server: &ws,
		Users: make(map[string]int),
		db: db,
	}
	ws.SetUpgrader(common.DefaultUpgrader)
	ws.SetEventHandlers(AdminWebSocket)
	ws.Handle(w, r)
}
