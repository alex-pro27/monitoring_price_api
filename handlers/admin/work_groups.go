package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Все рабочие группы пользователей
*/
func AllWorkGroups(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	var workGroups []models.WorkGroup
	db.Preload("Regions").Find(&workGroups, "active = true")
	var data []types.H
	for _, item := range workGroups {
		data = append(data, item.Serializer())
	}
	common.JSONResponse(w, data)
}
