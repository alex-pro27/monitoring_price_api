package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetRegions(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	type Region struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	regions := make([]*Region, 0)
	db.Model(&models.MonitoringGroups{}).Scan(&regions)
	common.JSONResponse(w, regions)
}
