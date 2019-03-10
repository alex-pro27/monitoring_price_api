package handlers

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Получить магазины для мониторинга(Конкуренты, маркетинг)
*/
func GetMonitoringShops(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*models.User)

	var monitoringShops []models.MonitoringShop

	var workGroupsIDX []uint
	for _, wg := range user.WorkGroup {
		workGroupsIDX = append(workGroupsIDX, wg.ID)
	}

	db := context.Get(r, "DB").(*gorm.DB)
	db.Preload(
		"Segments",
	).Joins(
		"INNER JOIN workgroup_monitoringshops wm ON wm.monitoring_shop_id = monitoring_shops.id",
	).Find(
		&monitoringShops, "wm.work_group_id IN (?)", workGroupsIDX,
	)

	var responseData []common.H

	for _, ms := range monitoringShops {
		responseData = append(responseData, ms.Serializer())
	}

	JSONResponse(w, responseData)
}
