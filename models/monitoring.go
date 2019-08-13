package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Monitoring struct {
	gorm.Model
	Name             string `gorm:"size:255;not null" form:"label:Название;required"`
	MonitoringTypeId uint
	MonitoringType   MonitoringType     `form:"label:Тип мониторинга"`
	MonitoringGroups []MonitoringGroups `gorm:"many2many:monitoring_groups_monitorings;" form:"label: Группы мониторинга"`
	MonitoringShops  []MonitoringShop   `gorm:"many2many:monitorings_monitoring_shops" form:"label:Магазины для мониторинга"`
	Users            []User             `gorm:"many2many:monitorings_users;" form:"label:Пользователи"`
	Active           bool               `gorm:"default:true" form:"label:Активная;type:switch"`
}

func (monitoring Monitoring) Serializer() types.H {
	var monitoringGroups []types.H
	for _, region := range monitoring.MonitoringGroups {
		monitoringGroups = append(monitoringGroups, region.Serializer())
	}
	return types.H{
		"id":                monitoring.ID,
		"name":              monitoring.Name,
		"monitoring_groups": monitoringGroups,
		"active":            monitoring.Active,
	}
}

func (Monitoring) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Мониторинг",
		Plural: "Мониторинги",
	}
}

func (Monitoring) Admin() types.AdminMeta {
	return types.AdminMeta{
		SearchFields: []string{"Name"},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "Active"},
		},
	}
}

func (monitoring Monitoring) String() string {
	return monitoring.Name
}
