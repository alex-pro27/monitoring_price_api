package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type WorkGroup struct {
	gorm.Model
	Name             string             `gorm:"size:255;not null" form:"label:Название;required"`
	Address          string             `gorm:"size:255" form:"label:Адрес"`
	MonitoringGroups []MonitoringGroups `gorm:"many2many:work_groups_monitoring_groups;" form:"label: Группы мониторинга"`
	MonitoringShops  []MonitoringShop   `gorm:"many2many:work_groups_monitoring_shops;" form:"label:Магазины для мониторинга"`
	Active           bool               `gorm:"default:true" form:"label:Активная;type:switch"`
}

func (workGroup WorkGroup) Serializer() types.H {
	var monitoringGroups []types.H
	for _, region := range workGroup.MonitoringGroups {
		monitoringGroups = append(monitoringGroups, region.Serializer())
	}
	return types.H{
		"id":                workGroup.ID,
		"name":              workGroup.Name,
		"address":           workGroup.Address,
		"monitoring_groups": monitoringGroups,
		"active":            workGroup.Active,
	}
}

func (WorkGroup) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Рабочая группа",
		Plural: "Рабочие группы",
	}
}

func (WorkGroup) Admin() types.AdminMeta {
	return types.AdminMeta{
		SearchFields: []string{"Name"},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "Address"},
			{Name: "Active"},
		},
	}
}

func (workGroup WorkGroup) String() string {
	return fmt.Sprintf("%s %s", workGroup.Name, workGroup.Address)
}
