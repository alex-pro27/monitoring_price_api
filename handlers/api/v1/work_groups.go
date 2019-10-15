package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetWorkGroups(w http.ResponseWriter, r *http.Request) {
	type WorkGroups struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		RegionID string `json:"region_id"`
	}

	data := make([]WorkGroups, 0)
	db := context.Get(r, "DB").(*gorm.DB)
	db.Model(&models.WorkGroup{}).Scan(&data)
	common.JSONResponse(w, data)
}
