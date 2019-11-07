package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetOnlineUsers(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	users := make([]models.User, 0)
	db.Find(&users, "online = true")
	data := make([]types.H, 0)
	for _, user := range users {
		data = append(data, user.Serializer())
	}
	common.JSONResponse(w, data)
}
