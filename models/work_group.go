package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type WorkGroup struct {
	gorm.Model
	Name            string           `gorm:"size:255;not null" form:"label:Название;required"`
	Address         string           `gorm:"size:255" form:"label:Адрес"`
	Regions         []Regions        `gorm:"many2many:work_groups_regions;" form:"label: Регионы"`
	MonitoringShops []MonitoringShop `gorm:"many2many:work_groups_monitoring_shops;" form:"label:Магазины для мониторинга"`
	Active          bool             `gorm:"default:true" form:"label:Активная;type:switch"`
}

func (workGroup WorkGroup) Serializer() types.H {
	var regions []types.H
	for _, region := range workGroup.Regions {
		regions = append(regions, region.Serializer())
	}
	return types.H{
		"id":      workGroup.ID,
		"name":    workGroup.Name,
		"address": workGroup.Address,
		"regions": regions,
		"active":  workGroup.Active,
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
