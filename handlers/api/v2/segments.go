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
Получить все доступные для юзера сегменты
*/
func GetSegments(w http.ResponseWriter, r *http.Request) {
	// TODO Добавить фильтр по периодам
	user := context.Get(r, "user").(*models.User)

	var monitoringIDX []uint
	for _, wg := range user.WorkGroups {
		for _, m := range wg.Monitorings {
			monitoringIDX = append(monitoringIDX, m.ID)
		}
	}
	var segments []models.Segment
	db := context.Get(r, "DB").(*gorm.DB)
	db.Preload(
		"Wares", "active = true",
	).Joins(
		"INNER JOIN monitoring_shops_segments ms ON ms.segment_id = segments.id",
	).Joins(
		"INNER JOIN work_groups_monitoring_shops wm ON wm.monitoring_shop_id = ms.monitoring_shop_id",
	).Select(
		"DISTINCT id, name, code, active",
	).Where(
		"wm.work_group_id IN (?) AND active = true", monitoringIDX,
	).Find(
		&segments,
	)

	var responseData []types.H

	for _, segment := range segments {
		if len(segment.Wares) > 0 {
			responseData = append(responseData, segment.Serializer())
		}
	}

	common.JSONResponse(w, responseData)
}
