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

	monitoringIDX := koazee.StreamOf(user.Monitorings).Map(func(m models.Monitoring) uint { return m.ID }).Out().Val()
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
		"INNER JOIN monitoring_shops_wares msw ON msw.ware_id = wares.id",
	).Joins(
		"INNER JOIN monitorings_monitoring_shops mms ON mms.monitoring_shop_id = msw.monitoring_shop_id",
	).Joins(
		"INNER JOIN monitorings m ON m.id = mms.monitoring_id",
	).Where("mms.monitoring_id IN (?)", monitoringIDX).Scan(&data)
	common.JSONResponse(w, data)
}
