package v1

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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

	var data []Ware

	db.Model(
		&models.Ware{},
	).Select(
		"DISTINCT wares.id, wares.code, wares.name, wares.barcode, wares.description, mt.id type_monitoring",
	).Joins(
		"INNER JOIN segments s ON s.id = wares.segment_id",
	).Joins(
		"INNER JOIN wares_monitoring_types wmt ON wmt.ware_id = wares.id",
	).Joins(
		"INNER JOIN monitoring_types mt ON wmt.monitoring_type_id = mt.id",
	).Joins(
		"INNER JOIN monitoring_shops_segments mss ON mss.segment_id = s.id",
	).Joins(
		"INNER JOIN monitoring_shops ms ON ms.id = mss.monitoring_shop_id",
	).Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON wgms.monitoring_shop_id = ms.id",
	).Joins(
		"INNER JOIN work_groups wg ON wg.id = wgms.work_group_id",
	).Joins(
		"INNER JOIN monitoring_types_periods mtp ON mtp.monitoring_type_id = mt.id",
	).Joins(
		"INNER JOIN periods p ON p.id = mtp.period_id",
	).Joins(
		"INNER JOIN work_groups_monitoring_groups wgmg ON wg.id = wgmg.work_group_id",
	).Joins(
		"INNER JOIN monitoring_groups mg ON mg.id = wgmg.monitoring_groups_id",
	).Where(
		"wares.active = true AND wg.name::text ~* ? AND mg.name::text ~* ? AND p.id IN (?)",
		vars["shop"],
		regions,
		periodsIDX,
	).Scan(&data)

	common.JSONResponse(w, data)

}
