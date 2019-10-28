package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"net/http"
)

func GetWares(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	user := context.Get(r, "user").(*models.User)
	periods := (&models.Period{}).Manager(db).GetAvailablePeriods()
	var periodsIDX []uint
	for _, period := range periods {
		periodsIDX = append(periodsIDX, period.ID)
	}

	type Ware struct {
		ID             uint   `json:"id"`
		Code           string `json:"code"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		Barcode        string `json:"barcode"`
		TypeMonitoring uint   `json:"type_monitoring"`
	}

	monitoringIDX := make([]uint, 0)
	for _, wg := range user.WorkGroups {
		for _, m := range wg.Monitorings {
			monitoringIDX = append(monitoringIDX, m.ID)
		}
	}
	monitoringIDX = koazee.StreamOf(monitoringIDX).RemoveDuplicates().Out().Val().([]uint)
	var data []Ware
	db.Model(new(models.Ware)).Select(
		"DISTINCT "+
			"wares.id, "+
			"wares.code, "+
			"wares.name, "+
			"wares.barcode, "+
			"wares.description, "+
			"m.monitoring_type_id type_monitoring",
	).Joins(
		"INNER JOIN monitorings_wares mw ON mw.ware_id = wares.id",
	).Joins(
		"INNER JOIN monitorings m ON m.id = mw.monitoring_id",
	).Joins(
		"INNER JOIN monitoring_types mt ON mt.id = m.monitoring_type_id",
	).Joins(
		"INNER JOIN monitoring_types_periods mtp ON mtp.monitoring_type_id = mt.id",
	).Joins(
		"INNER JOIN periods p ON p.id = mtp.period_id",
	).Where(
		"mw.monitoring_id IN (?) AND p.id IN (?)",
		monitoringIDX,
		periodsIDX,
	).Scan(&data)
	common.JSONResponse(w, data)
}
