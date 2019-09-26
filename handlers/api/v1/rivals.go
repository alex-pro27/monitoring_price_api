package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

/**
Конкуренты => сегметы => IDs товаров
*/
func GetRivals(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	user := context.Get(r, "user").(*models.User)
	periods := (new(models.Period)).Manager(db).GetAvailablePeriods()
	var periodsIDX []uint
	for _, period := range periods {
		periodsIDX = append(periodsIDX, period.ID)
	}
	workGroupIDX := make([]uint, 0)
	monitoringIDX := make([]uint, 0)
	for _, wg := range user.WorkGroups {
		workGroupIDX = append(workGroupIDX, wg.ID)
		for _, m := range wg.Monitorings {
			monitoringIDX = append(monitoringIDX, m.ID)
		}
	}
	var rivals []models.MonitoringShop
	db.Preload("Segments").Select(
		"DISTINCT monitoring_shops.*",
	).Joins(
		"INNER JOIN work_groups_monitoring_shops wgms ON wgms.monitoring_shop_id = monitoring_shops.id",
	).Find(
		&rivals,
		"monitoring_shops.active = true AND wgms.work_group_id IN (?)",
		workGroupIDX,
	)

	var wares []struct {
		WareID    uint
		SegmentID uint
	}

	db.Model(models.Ware{}).Select("wares.id ware_id, segment_id").Joins(
		"INNER JOIN monitorings_wares mw ON mw.ware_id = wares.id",
	).Where("mw.monitoring_id IN (?)", monitoringIDX).Scan(&wares)

	_segments := make(map[uint]models.Segment)
	waresIDxByRivalIDBySegmentID := make(map[uint]map[uint][]uint)

	for _, rival := range rivals {
		if waresIDxByRivalIDBySegmentID[rival.ID] == nil {
			waresIDxByRivalIDBySegmentID[rival.ID] = map[uint][]uint{}
		}
		for _, s := range rival.Segments {
			for _, ws := range wares {
				if ws.SegmentID == s.ID {
					if waresIDxByRivalIDBySegmentID[rival.ID][ws.SegmentID] == nil {
						waresIDxByRivalIDBySegmentID[rival.ID][ws.SegmentID] = make([]uint, 0)
					}
					waresIDxByRivalIDBySegmentID[rival.ID][ws.SegmentID] = append(waresIDxByRivalIDBySegmentID[rival.ID][ws.SegmentID], ws.WareID)
					_segments[ws.SegmentID] = s
				}
			}
		}
	}

	var data []types.H

	for _, rival := range rivals {
		var segments []types.H
		for segmentID, waresIDX := range waresIDxByRivalIDBySegmentID[rival.ID] {
			if len(waresIDX) == 0 {
				continue
			}

			_segment := _segments[segmentID]

			segments = append(segments, types.H{
				"id":    _segment.ID,
				"code":  _segment.Code,
				"name":  _segment.Name,
				"wares": waresIDX,
			})
		}
		data = append(data, types.H{
			"id":            rival.ID,
			"name":          rival.Name,
			"address":       rival.Address,
			"segments":      segments,
			"is_must_photo": rival.IsMustPhoto,
		})
	}

	common.JSONResponse(w, data)
}
