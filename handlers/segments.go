package handlers

import (
	. "github.com/alex-pro27/monitoring_price_api/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Получить все доступный сегменты
*/
func GetSegments(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*models.User)

	var workGroupsIDX []uint
	for _, wg := range user.WorkGroup {
		workGroupsIDX = append(workGroupsIDX, wg.ID)
	}
	var segments []models.Segment
	db := context.Get(r, "DB").(*gorm.DB)
	db.Preload(
		"Wares",
	).Joins(
		"INNER JOIN monitoringshops_segments ms ON ms.segment_id = segments.id",
	).Joins(
		"INNER JOIN workgroup_monitoringshops wm ON wm.monitoring_shop_id = ms.monitoring_shop_id",
	).Select(
		"DISTINCT id, name, code",
	).Where(
		"wm.work_group_id IN (?)", workGroupsIDX,
	).Find(
		&segments,
	)

	var responseData []H

	for _, segment := range segments {
		responseData = append(responseData, segment.Serializer())
	}

	JSONResponse(w, responseData)
}
