package v1

import (
	"github.com/alex-pro27/monitoring_price_api/handlers/common"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/models"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetPeriods(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "DB").(*gorm.DB)
	var monitoringTypes []models.MonitoringType
	db.Preload("Periods").Find(&monitoringTypes)
	var data []types.H
	for _, item := range monitoringTypes {
		if len(item.Periods) > 0 {
			periodDates := item.Periods[0].GetPeriodDates()
			start := periodDates.Dates[0].Format(helpers.ISO8601)
			end := periodDates.Dates[1].Format(helpers.ISO8601)
			data = append(data, types.H{
				"id":              item.ID,
				"period_name":     item.Periods[0].GetPeriodName(),
				"monitoring_name": item.Name,
				"start":           start,
				"to":              end,
				"end":             end,
			})
		}
	}

	common.JSONResponse(w, data)
}
