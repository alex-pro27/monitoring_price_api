package v1

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"net/http"
	"strings"
)

func GetWares(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	periods := (&models.Period{}).Manager(db).GetAvailablePeriods()
	vars := mux.Vars(r)
	_regions := strings.Builder{}
	for _, region := range strings.Split(vars["region"], "-") {
		_regions.Write([]byte(fmt.Sprintf("(%s)|", region)))
	}
	regions := _regions.String()
	regions = regions[:len(regions)-1]

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

	var monitorings []models.Monitoring
	db.Select(
		"DISTINCT monitorings.*",
	).Joins(
		"INNER JOIN monitoring_types mt ON mt.id = monitoring_type_id",
	).Joins(
		"INNER JOIN work_groups_monitorings wgm ON wgm.monitoring_id = monitorings.id",
	).Joins(
		"INNER JOIN work_groups wg ON wg.id = wgm.work_group_id",
	).Joins(
		"INNER JOIN monitoring_groups mg ON mg.id = monitorings.monitoring_group_id",
	).Joins(
		"INNER JOIN monitoring_types_periods mtp ON mtp.monitoring_type_id = mt.id",
	).Where(
		"monitorings.active = true AND mtp.period_id IN (?) AND wg.name ilike ? AND mg.name ilike ?",
		periodsIDX,
		vars["shop"],
		vars["region"],
	).Find(&monitorings)
	monitoringIDX := koazee.StreamOf(monitorings).Map(func(m models.Monitoring) uint { return m.ID }).Out().Val()

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
	).Where("mw.monitoring_id IN (?)", monitoringIDX).Scan(&data)
	common.JSONResponse(w, data)
}
