package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Типы мониторинга
*/
type MonitoringType struct {
	gorm.Model
	Name    string   `gorm:"size:255;not null" form:"label:Имя"`
	Periods []Period `gorm:"many2many:monitoring_types_periods" form:"label:Перидоы"`
	Active  bool     `gorm:"default:true" form:"label:Активный;type:switch"`
}

func (MonitoringType) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Тип мониторинга",
		Plural: "Типы мониторинга",
	}
}

func (monitoringType MonitoringType) String() string {
	return monitoringType.Name
}
