package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Группы мониторинга
*/
type MonitoringGroups struct {
	gorm.Model
	Name       	string      `gorm:"size:255" form:"label:Название"`
	WorkGroups 	[]WorkGroup `gorm:"many2many:work_groups_monitoring_groups;" form:"label:Рабочие группы"`
	Wares 		[]Ware		`gorm:"many2many:wares_monitoring_groups" form:"label:Товары"`
}

func (region MonitoringGroups) Serializer() types.H {
	return types.H{
		"id":   region.ID,
		"name": region.Name,
	}
}

func (MonitoringGroups) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Группа мониторинга",
		Plural: "Группы мониторинга",
	}
}

func (region MonitoringGroups) String() string {
	return region.Name
}
