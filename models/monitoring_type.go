package models

import "github.com/jinzhu/gorm"

/**
Типы мониторинга
*/
type MonitoringType struct {
	gorm.Model
	Name    string   `gorm:"size:255;not null"`
	Periods []Period `gorm:"many2many:monitoring_types_periods"`
}
