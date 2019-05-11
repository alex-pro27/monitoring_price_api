package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetMonitoringShops(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	type MonitoringShop struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Code    string `json:"code"`
		Address string `json:"address"`
	}
	monitoringShops := make([]*MonitoringShop, 0)
	db.Model(&models.MonitoringShop{}).Scan(&monitoringShops)
	common.JSONResponse(w, monitoringShops)
}
