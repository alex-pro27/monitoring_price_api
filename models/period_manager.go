package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"time"
)

type PeriodManager struct {
	*gorm.DB
	self *Period
}

func (manager *PeriodManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.self.Period = PeriodsType(fields["period"].(float64))
	manager.DB.Create(manager.self)
	return nil
}

func (manager *PeriodManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.self.Period = PeriodsType(fields["period"].(float64))
	manager.Save(manager.self)
	return nil
}

func (manager *PeriodManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}

func (manager *PeriodManager) GetAvailablePeriods() []Period {
	var all, periods []Period
	currentDate := time.Now()
	manager.Find(&all)
	for _, period := range all {
		periodDates := period.GetPeriodDates()
		if periodDates.Type == PERIOD_DAY {
			for _, date := range periodDates.Dates {
				if date.Day() == currentDate.Day() {
					periods = append(periods, period)
					break
				}
			}
		} else {
			if currentDate.Unix() > periodDates.Dates[0].Unix() && currentDate.Unix() < periodDates.Dates[1].Unix() {
				periods = append(periods, period)
			}
		}
	}
	return periods
}
