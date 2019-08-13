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
	monitoringIDX := make([]uint, 0)

	for _, m := range user.Monitorings {
		monitoringIDX = append(monitoringIDX, m.ID)
	}
	var rivals []models.MonitoringShop
	db.Preload("Wares.Segment").Select(
		"DISTINCT monitoring_shops.*",
	).Joins(
		"INNER JOIN monitorings_monitoring_shops mms ON mms.monitoring_shop_id = monitoring_shops.id",
	).Find(
		&rivals,
		"monitoring_shops.active = true AND mms.monitoring_id IN (?)",
		monitoringIDX,
	)

	_segments := make(map[uint]models.Segment)
	waresIDxByRivalIDBySegmentID := make(map[uint]map[uint][]uint)

	for _, rival := range rivals {
		if waresIDxByRivalIDBySegmentID[rival.ID] == nil {
			waresIDxByRivalIDBySegmentID[rival.ID] = map[uint][]uint{}
		}
		for _, w := range rival.Wares {
			if waresIDxByRivalIDBySegmentID[rival.ID][w.SegmentId] == nil {
				waresIDxByRivalIDBySegmentID[rival.ID][w.SegmentId] = make([]uint, 0)
			}
			waresIDxByRivalIDBySegmentID[rival.ID][w.SegmentId] = append(waresIDxByRivalIDBySegmentID[rival.ID][w.SegmentId], w.ID)
			_segments[w.SegmentId] = w.Segment
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
