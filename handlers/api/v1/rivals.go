package v1

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

/**
Конкуренты => сегметы => IDs товаров
*/
func GetRivals(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	periods := (&models.Period{}).Manager(db).GetAvailablePeriods()
	vars := mux.Vars(r)
	_regions := strings.Builder{}
	for _, region := range strings.Split(vars["region"], "-") {
		_regions.Write([]byte(fmt.Sprintf("(%s)|", region)))
	}
	regions := _regions.String()
	regions = regions[:len(regions)-1]
	var rivals []models.MonitoringShop

	var periodsIDX []uint
	for _, period := range periods {
		periodsIDX = append(periodsIDX, period.ID)
	}

	db.Preload(
		"Segments.Wares",
	).Preload(
		"WorkGroup.Regions",
	).Select(
		"DISTINCT monitoring_shops.*",
	).Joins(
		"INNER JOIN monitoring_shops_segments mss ON mss.monitoring_shop_id = monitoring_shops.id",
	).Joins(
		"INNER JOIN segments s ON mss.segment_id = s.id",
	).Joins(
		"INNER JOIN wares w ON w.segment_id = s.id",
	).Joins(
		"INNER JOIN wares_monitoring_types wmt ON wmt.ware_id = w.id",
	).Joins(
		"INNER JOIN monitoring_types mt ON mt.id = wmt.monitoring_type_id",
	).Joins(
		"INNER JOIN monitoring_types_periods mtp ON mtp.monitoring_type_id = mt.id",
	).Joins(
		"INNER JOIN periods p ON p.id = mtp.period_id",
	).Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON monitoring_shops.id = wgms.monitoring_shop_id",
	).Joins(
		"INNER JOIN work_groups wg ON wg.id = wgms.work_group_id",
	).Joins(
		"INNER JOIN work_groups_monitoring_groups wgmg ON wg.id = wgmg.work_group_id",
	).Joins(
		"INNER JOIN monitoring_groups mg ON mg.id = wgmg.monitoring_groups_id",
	).Find(
		&rivals,
		"monitoring_shops.active = true AND wg.name::text ~* ? AND mg.name::text ~* ? AND p.id IN (?)",
		vars["shop"],
		regions,
		periodsIDX,
	)

	var data []types.H

	for _, rival := range rivals {
		var segments []types.H
		for _, segment := range rival.Segments {
			if len(segment.Wares) == 0 {
				continue
			}
			var waresIDX []uint
			for _, ware := range segment.Wares {
				waresIDX = append(waresIDX, ware.ID)
			}

			segments = append(segments, types.H{
				"id":    segment.ID,
				"code":  segment.Code,
				"name":  segment.Name,
				"wares": waresIDX,
			})
		}
		data = append(data, types.H{
			"id":            rival.ID,
			"name":          rival.Name,
			"address":       rival.Address,
			"region":        rival.WorkGroup[0].MonitoringGroups[0].Name,
			"segments":      segments,
			"is_must_photo": rival.IsMustPhoto,
		})
	}

	common.JSONResponse(w, data)
}
