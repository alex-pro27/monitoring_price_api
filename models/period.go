package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/lib/pq"
	"time"
)

type PeriodsType int

const (
	PERIOD_DAY PeriodsType = iota
	PERIOD_WEEK
	PERIOD_MONTH
	PERIOD_QUARTER
	PERIOD_YEAR
)

var PeriodChoices = map[PeriodsType]string{
	PERIOD_DAY:     "День",
	PERIOD_WEEK:    "Неделя",
	PERIOD_MONTH:   "Месяц",
	PERIOD_QUARTER: "Квартал",
	PERIOD_YEAR:    "Год",
}

type CurrentPeriods struct {
	Type  PeriodsType `json:"type"`
	Dates []time.Time `json:"dates"`
}

type Period struct {
	gorm.Model
	Period       PeriodsType `form:"choice:GetPeriodChoices"`
	Start        int
	End          int
	SelectedDays pq.Int64Array `gorm:"type:integer[]"`
}

func (Period) GetPeriodChoices() map[PeriodsType]string {
	return PeriodChoices
}

func (period Period) GetPeriodName() string {
	return PeriodChoices[period.Period]
}

func (period *Period) CRUD(db *gorm.DB) types.CRUDManager {
	return period.Manager(db)
}

func (period *Period) Manager(db *gorm.DB) *PeriodManager {
	return &PeriodManager{db, period}
}

func (Period) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Период",
		Plural: "Периоды",
	}
}

func (period Period) String() string {
	return period.GetPeriodName()
}

func (period Period) GetPeriodDates() (currentPeriods CurrentPeriods) {
	year, month, day := time.Now().Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	dayDuration := time.Hour * 24
	switch period.Period {
	case PERIOD_DAY:
		currentPeriods.Type = PERIOD_DAY
		for _, _day := range period.SelectedDays {
			currentPeriods.Dates = append(
				currentPeriods.Dates,
				time.Date(year, month, int(_day), 0, 0, 0, 0, time.Local),
			)
		}
		break
	case PERIOD_WEEK:
		currentPeriods.Type = PERIOD_WEEK
		for date.Weekday() != time.Monday {
			date = date.AddDate(0, 0, -1)
		}
		end := period.End
		if end <= period.Start {
			end += 7
		}
		currentPeriods.Dates = []time.Time{
			date.AddDate(0, 0, period.Start-1),
			date.AddDate(0, 0, end-1),
		}
		break
	case PERIOD_MONTH:
		currentPeriods.Type = PERIOD_MONTH
		var start, end time.Time
		beginningMonth := now.BeginningOfMonth()
		prevMonth := beginningMonth.Add(-(dayDuration))
		prevMonth = time.Date(prevMonth.Year(), prevMonth.Month(), period.Start, 0, 0, 0, 0, time.Local)
		tempDateEnd := beginningMonth.Add(dayDuration*time.Duration(period.End) - 1)
		if day > period.End {
			// Если текущий день больше дня окончания мониторинга то выбираем следующий месяц
			start = prevMonth.Add(dayDuration * time.Duration(period.Start))
			end = start.Add(dayDuration * time.Duration(period.End))
		} else if tempDateEnd.Sub(prevMonth).Hours()/24 <= float64(now.EndOfMonth().Day()) {
			// Если разница в кол. пройденых дней не превышает или равно кол. дней в пред. месяце
			// то выбираем предыдущий месяц
			start = prevMonth
			end = time.Date(tempDateEnd.Year(), tempDateEnd.Month(), tempDateEnd.Day(), 0, 0, 0, 0, time.Local)
		} else {
			// Иначе остаемся в текущем месяц
			start = time.Date(year, month, period.Start, 0, 0, 0, 0, time.Local)
			end = time.Date(year, month, period.End, 0, 0, 0, 0, time.Local)
		}
		currentPeriods.Dates = []time.Time{start, end}
		break
	case PERIOD_QUARTER, PERIOD_YEAR:
		if period.Period == PERIOD_YEAR {
			currentPeriods.Type = PERIOD_YEAR
		} else {
			currentPeriods.Type = PERIOD_QUARTER
		}
		startYear := now.BeginningOfYear()
		currentPeriods.Dates = []time.Time{
			startYear.Add(dayDuration * time.Duration(period.Start)),
			startYear.Add(dayDuration * time.Duration(period.End)),
		}
		break
	}
	return currentPeriods
}
