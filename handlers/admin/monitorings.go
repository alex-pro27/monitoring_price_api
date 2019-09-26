package admin

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetAllMonitoringList(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	monitorings := make([]models.Monitoring, 0)
	db.Find(&monitorings)
	data := make([]types.H, 0)
	for _, it := range monitorings {
		data = append(data, it.Serializer())
	}
	common.JSONResponse(w, data)
}
