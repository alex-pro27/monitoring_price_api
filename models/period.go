package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
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

type Period struct {
	gorm.Model
	Period       PeriodsType `form:"choice:GetPeriodChoices"`
	Start        time.Time
	End          time.Time
	SelectedDays pq.Int64Array `gorm:"type:integer[]"`
}

func (Period) GetPeriodChoices() map[PeriodsType]string {
	return PeriodChoices
}

func (period Period) GetPeriodName() string {
	return PeriodChoices[period.Period]
}

func (period *Period) CRUD(db *gorm.DB) types.CRUDManager {
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
