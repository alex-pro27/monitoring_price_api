package admin

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/databases"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"net/http"
)

type AdminWebSocketHandler struct {
	Server *common.WebSocket
	Users  map[string][]int
	IsInit bool
	Decorators []func(handleFunc common.WSHandleFunc) common.WSHandleFunc
	db *gorm.DB
}

func (h *AdminWebSocketHandler) Init(server *common.WebSocket)  {
	h.Server = server
	h.Server.SetUpgrader(common.DefaultUpgrader)
	h.Server.SetEventHandlers(h)
	h.Users = make(map[string][]int)
	h.IsInit = true
	h.Decorators = []func(handleFunc common.WSHandleFunc) common.WSHandleFunc {
		/**Handle error*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H) {
				if rec := recover(); rec != nil {
					logger.Logger.Errorf("websocket error: %v", rec)
				}
				logger.Logger.Infof("websocket: client send message")
				f(clientID, message)
			}
		},
		/**Connect db*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H) {
				h.db = databases.ConnectDefaultDB()
				f(clientID, message)
				defer logger.HandleError(h.db.Close())
			}
		},
		/**Check user*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H) {
				user := models.User{}
				token := message["token"]
				if token != nil {
					user.Manager(h.db).GetUserByToken(message["token"].(string))
				}
				if !user.IsStaff {
					h.Server.Emit(clientID, "on_connect", types.H{
						"error":   true,
						"code":    403,
						"message": "Permission denied",
					})
					logger.Logger.Warning("Admin Websocket error: Permission denied")
					client := h.Server.Client(clientID)
					if client != nil {
						logger.HandleError(client.Close())
					}
					return
				}
				if h.Users[user.Token.Key] == nil {
					h.Users[user.Token.Key] = make([]int, 0)
				}
				h.Users[user.Token.Key] = append(h.Users[user.Token.Key], clientID)
				f(clientID, message)
			}
		},
	}
}

func (h AdminWebSocketHandler) Emit(token, event string, message types.H) {
	clientIDX := AdminWebSocket.Users[token]
	for _, clientID := range clientIDX {
		h.Server.Emit(clientID, event, message)
	}
}

func (h *AdminWebSocketHandler) OnOpen(clientID int) {
	h.Server.Emit(clientID, "on_open", types.H{
		"message": fmt.Sprintf("client on connected %d", clientID),
	})
}

func (h *AdminWebSocketHandler) OnClose(clientID int) {
	for token, idx := range h.Users {
		if len(idx) > 0 {
			for i, id := range idx {
				if id == clientID {
					h.Users[token] = append(h.Users[token][:i], h.Users[token][i+1:]...)
					break
				}
			}
		}
	}
}

func (h *AdminWebSocketHandler) OnConnect(clientID int, message types.H) {
	h.Server.Emit(clientID, "on_connect", types.H{
		"message": "Connected!",
	})
}

var AdminWebSocket = &AdminWebSocketHandler{}

func StartWebsocket(w http.ResponseWriter, r *http.Request) {
	if !AdminWebSocket.IsInit {
		AdminWebSocket.Init(&common.WebSocket{})
	}
	AdminWebSocket.Server.Handle(w, r)
}
