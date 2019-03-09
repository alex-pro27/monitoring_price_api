package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

/**
Регион проведения мониторинга
*/
type Regions struct {
	gorm.Model
	Name       string      `gorm:"size:255"`
	WorkGroups []WorkGroup `gorm:"many2many:workgroup_regions;"`
}

func (region Regions) Serializer() common.H {
	return common.H{
		"id":   region.ID,
		"name": region.Name,
	}
}
