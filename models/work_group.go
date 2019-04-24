package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type WorkGroup struct {
	gorm.Model
	Name            string           `gorm:"size:255;not null"`
	Address         string           `gorm:"size:255"`
	Regions         []Regions        `gorm:"many2many:work_groups_regions;"`
	MonitoringShops []MonitoringShop `gorm:"many2many:work_groups_monitoring_shops;"`
	Active          bool             `gorm:"default:true"`
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

func (workGroup WorkGroup) String() string {
	return fmt.Sprintf("%s %s", workGroup.Name, workGroup.Address)
}
