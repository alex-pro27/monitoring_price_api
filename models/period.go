package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type Periods int

const (
	DAY     Periods = 0
	WEEK    Periods = 1
	MONTH   Periods = 2
	QUARTER Periods = 3
	YEAR    Periods = 4
)

var PeriodChoices = map[Periods]string{
	DAY:     "День",
	WEEK:    "Неделя",
	MONTH:   "Месяц",
	QUARTER: "Квартал",
	YEAR:    "Год",
}

type Period struct {
	gorm.Model
	Period       Periods
	Start        time.Time
	End          time.Time
	SelectedDays pq.Int64Array `gorm:"type:integer[]"`
}

func (period Period) GetPeriodName() string {
	return PeriodChoices[period.Period]
}
