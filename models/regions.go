package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Регион
*/
type Region struct {
	gorm.Model
	Name        string       `gorm:"size:255" form:"label:Название"`
	Monitorings []Monitoring `gorm:"many2many:regions_monitorings;" form:"label:Мониториги"`
}

func (Region) Admin() types.AdminMeta {
	return types.AdminMeta{
		OrderBy: []string{"Name"},
	}
}

func (region Region) Serializer() types.H {
	return types.H{
		"id":   region.ID,
		"name": region.Name,
	}
}

func (Region) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Регион",
		Plural: "Регионы",
	}
}

func (region Region) String() string {
	return region.Name
}
