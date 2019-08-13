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
	Name        string       `gorm:"size:255" form:"label:Название"`
	Monotorings []Monitoring `gorm:"many2many:monitoring_groups_monitorings;" form:"label:Мониториги"`
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
