package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/wesovilabs/koazee"
	"net/http"
)

/**
Конкуренты => сегметы => IDs товаров
*/
func GetRivals(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	periods := (new(models.Period)).Manager(db).GetAvailablePeriods()
	vars := mux.Vars(r)

	var periodsIDX []uint
	for _, period := range periods {
		periodsIDX = append(periodsIDX, period.ID)
	}

	var monitorings []models.Monitoring
	db.Preload("WorkGroups", func(db *gorm.DB) *gorm.DB {
		return db.Where("name ilike ?", vars["shop"])
	}).Select(
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

	monitoringIDX := make([]uint, 0)
	workGroupIDX := make([]uint, 0)

	for _, m := range monitorings {
		monitoringIDX = append(monitoringIDX, m.ID)
		for _, wg := range m.WorkGroups {
			workGroupIDX = append(workGroupIDX, wg.ID)
		}
	}

	workGroupIDX = koazee.StreamOf(workGroupIDX).RemoveDuplicates().Out().Val().([]uint)

	var rivals []models.MonitoringShop
	db.Select(
		"DISTINCT monitoring_shops.*",
	).Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON wgms.monitoring_shop_id = monitoring_shops.id",
	).Find(
		&rivals,
		"monitoring_shops.active = true AND wgms.work_group_id IN (?)",
		workGroupIDX,
	)

	var _segments []models.Segment
	db.Select(
		"DISTINCT segments.*",
	).Preload(
		"Wares", func(db *gorm.DB) *gorm.DB {
			return db.Joins(
				"INNER JOIN monitorings_wares mw ON mw.ware_id = wares.id",
			).Where("mw.monitoring_id IN (?)", monitoringIDX)
		},
	).Joins(
		"INNER JOIN wares w ON w.segment_id = segments.id",
	).Joins(
		"INNER JOIN monitorings_wares mw ON mw.ware_id = w.id",
	).Find(&_segments, "active = true AND mw.monitoring_id IN (?)", monitoringIDX)

	_monitoringGroup := new(models.MonitoringGroups)

	db.First(&_monitoringGroup, "name::text ~* ?", vars["region"])

	var data []types.H

	for _, rival := range rivals {
		var segments []types.H
		for _, segment := range _segments {
			var waresIDX []uint
			for _, ware := range segment.Wares {
				if ware.SegmentId == segment.ID {
					waresIDX = append(waresIDX, ware.ID)
				}
			}

			if len(waresIDX) == 0 {
				continue
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
			"region":        _monitoringGroup.Name,
			"segments":      segments,
			"is_must_photo": rival.IsMustPhoto,
		})
	}

	common.JSONResponse(w, data)
}
