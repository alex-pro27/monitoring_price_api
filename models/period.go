package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/wesovilabs/koazee"
	"regexp"
	"strconv"
	"strings"
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
	Name         string        `gorm:"type:varchar(255)" form:"label:Название"`
	Period       PeriodsType   `form:"choice:GetPeriodChoices;label:Период"`
	Start        int           `form:"label:Начало периода"`
	End          int           `form:"label:Конец периода"`
	SelectedDays pq.Int64Array `gorm:"type:integer[]" form:"label:Дни"`
	Active       bool          `gorm:"default:true" form:"label:Активный;type:switch"`
}

func (Period) GetPeriodChoices() map[PeriodsType]string {
	return PeriodChoices
}

func (period Period) GetPeriodName() string {
	return PeriodChoices[period.Period]
}

func (period *Period) SetSelectedDays(days string) (err error) {
	_days := make(pq.Int64Array, 0)
	if days != "" {
		if match, _ := regexp.MatchString("^[1-7],\\s?([1-7],\\s?){0,5}[1-7]?$", days); !match {
			return errors.New("Введите дни от 1 до 7 через запятую")
		}
		for _, day := range strings.Split(days, ",") {
			_day, _ := strconv.ParseInt(day, 10, 64)
			if _day != 0 {
				_days = append(_days, _day)
			}
		}
		_days = koazee.StreamOf(_days).RemoveDuplicates().Out().Val().([]int64)
	}
	period.SelectedDays = _days
	return nil
}

func (period *Period) CRUD(db *gorm.DB) types.CRUDManager {
	return period.Manager(db)
}

func (period Period) GetCurrentPeriod() string {
	var dates []string
	for _, date := range period.GetPeriodDates().Dates {
		dates = append(dates, date.Format(helpers.ISO8601))
	}
	return strings.Join(dates, ",")
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

func (period *Period) Admin() types.AdminMeta {
	return types.AdminMeta{
		SortFields: []types.AdminMetaField{
			{
				Name: "Name",
			},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "GetPeriodName",
				Label: "Период",
			},
			{
				Name:   "GetCurrentPeriod",
				Label:  "Даты",
				ToHTML: "date",
			},
		},
	}
}

func (period Period) String() string {
	if period.Name == "" {
		return period.GetPeriodName()
	} else {
		return period.Name
	}
}

func (period Period) GetPeriodDates() (currentPeriods CurrentPeriods) {
	year, month, day := time.Now().Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	dayDuration := time.Hour * 24
	switch period.Period {
	case PERIOD_DAY:
		currentPeriods.Type = PERIOD_DAY
		for _, _day := range period.SelectedDays {
			currentPeriods.Dates = append(
				currentPeriods.Dates,
				date.AddDate(0, 0, int(_day-1)),
			)
		}
		break
	case PERIOD_WEEK:
		currentPeriods.Type = PERIOD_WEEK
		end := period.End
		start := period.Start
		if end <= start {
			if time.Weekday(start) > time.Now().Weekday() {
				start -= 7
			} else {
				end += 7
			}
		}
		currentPeriods.Dates = []time.Time{
			date.AddDate(0, 0, start-1),
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
			start = beginningMonth.Add(dayDuration * time.Duration(period.Start-1))
			end = now.EndOfMonth().Add(dayDuration * time.Duration(period.End))
		} else if tempDateEnd.Sub(prevMonth).Hours() / 24 <= float64(now.EndOfMonth().Day()) || day <= period.End && period.End - period.Start < 8 {
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
	case PERIOD_QUARTER:
		currentPeriods.Type = PERIOD_QUARTER
		startQuarter := now.BeginningOfQuarter()
		currentPeriods.Dates = []time.Time{
			startQuarter.Add(dayDuration * time.Duration(period.Start-1)),
			startQuarter.Add(dayDuration * time.Duration(period.End-1)),
		}
	case PERIOD_YEAR:
		if period.Period == PERIOD_YEAR {
			currentPeriods.Type = PERIOD_YEAR
		}
		startYear := now.BeginningOfYear()
		currentPeriods.Dates = []time.Time{
			startYear.Add(dayDuration * time.Duration(period.Start-1)),
			startYear.Add(dayDuration * time.Duration(period.End-1)),
		}
		break
	}
	return currentPeriods
}
