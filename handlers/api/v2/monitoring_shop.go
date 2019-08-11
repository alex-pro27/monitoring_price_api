package v2

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Получить магазины для мониторинга(Конкуренты, маркетинг)
*/
func GetMonitoringShops(w http.ResponseWriter, r *http.Request) {
	// TODO Добавить фильтр по периодам
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
		"INNER JOIN work_groups_monitoring_shops wm ON wm.monitoring_shop_id = monitoring_shops.id",
	).Find(
		&monitoringShops, "wm.work_group_id IN (?)", workGroupsIDX,
	)

	var responseData []types.H

	//for _, ms := range monitoringShops {
	//	if len(ms.Wares) > 0 {
	//		responseData = append(responseData, ms.Serializer())
	//	}
	//}

	common.JSONResponse(w, responseData)
}
