package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Регион проведения мониторинга
*/
type Regions struct {
	gorm.Model
	Name       string      `gorm:"size:255"`
	WorkGroups []WorkGroup `gorm:"many2many:work_groups_regions;"`
}

func (region Regions) Serializer() types.H {
	return types.H{
		"id":   region.ID,
		"name": region.Name,
	}
}

func (Regions) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Регион",
		Plural: "Регионы",
	}
}

func (region Regions) String() string {
	return region.Name
}
