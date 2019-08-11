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
	Server     *common.WebSocket
	Users      map[string][]int
	IsInit     bool
	Decorators []func(handleFunc common.WSHandleFunc) common.WSHandleFunc
	db         *gorm.DB
}

func (h *AdminWebSocketHandler) Init(server *common.WebSocket) {
	h.Server = server
	h.Server.SetUpgrader(common.DefaultUpgrader)
	h.Server.SetEventHandlers(h)
	h.Users = make(map[string][]int)
	h.IsInit = true
	h.Decorators = []func(handleFunc common.WSHandleFunc) common.WSHandleFunc{
		/**Handle error*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H, args ...interface{}) {
				if rec := recover(); rec != nil {
					logger.Logger.Errorf("websocket error: %v", rec)
				}
				logger.Logger.Infof("websocket: client send message")
				f(clientID, message)
			}
		},
		/**Connect db*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H, args ...interface{}) {
				h.db = databases.ConnectDefaultDB()
				f(clientID, message)
				logger.HandleError(h.db.Close())
			}
		},
		/**Check user*/
		func(f common.WSHandleFunc) common.WSHandleFunc {
			return func(clientID int, message types.H, args ...interface{}) {
				user := new(models.User)
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
				f(clientID, message, user)
			}
		},
	}
}

func (h AdminWebSocketHandler) Emit(token, event string, message types.H) {
	if h.IsInit {
		clientIDX := AdminWebSocket.Users[token]
		for _, clientID := range clientIDX {
			h.Server.Emit(clientID, event, message)
		}
	}
}

func (h AdminWebSocketHandler) EmitAll(excludeClientID int, event string, message types.H) {
	if h.IsInit {
		for _, clientIDX := range h.Users {
			for _, clientID := range clientIDX {
				if excludeClientID != clientID {
					h.Server.Emit(clientID, event, message)
				}
			}
		}
	}
}

func (h *AdminWebSocketHandler) OnOpen(clientID int) {
	h.Server.Emit(clientID, "on_open", types.H{
		"message": fmt.Sprintf("client on connected %d", clientID),
	})
}

func (h *AdminWebSocketHandler) OnClose(clientID int) {
	var token string
	for _token, idx := range h.Users {
		if len(idx) > 0 {
			for i, id := range idx {
				if id == clientID {
					h.Users[_token] = append(h.Users[_token][:i], h.Users[_token][i+1:]...)
					if len(h.Users[_token]) == 0 {
						token = _token
					}
					break
				}
			}
		}
	}
	if token != "" {
		db := databases.ConnectDefaultDB()
		user := new(models.User)
		user.Manager(db).GetUserByToken(token)
		user.Online = false
		db.Save(user)
		logger.HandleError(h.db.Close())
		h.EmitAll(-1, "client_leaved", map[string]interface{}{
			"client_name": user.GetFullName(),
		})
	}
}

func (h *AdminWebSocketHandler) OnConnect(clientID int, message types.H, args ...interface{}) {
	user := args[0].(*models.User)
	user.Online = true
	h.db.Save(user)
	h.Server.Emit(clientID, "on_connect", types.H{
		"message": "Connected!",
	})
	if len(h.Users[user.Token.Key]) == 1 {
		h.EmitAll(clientID, "client_joined", map[string]interface{}{
			"client_name": user.GetFullName(),
		})
	}
}

var AdminWebSocket = new(AdminWebSocketHandler)

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	if !AdminWebSocket.IsInit {
		AdminWebSocket.Init(new(common.WebSocket))
	}
	AdminWebSocket.Server.Handle(w, r)
}
