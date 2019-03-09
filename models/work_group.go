package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

type WorkGroup struct {
	gorm.Model
	Name            string           `gorm:"size:255;not null"`
	Address         string           `gorm:"size:255"`
	Regions         []Regions        `gorm:"many2many:workgroup_regions;"`
	MonitoringShops []MonitoringShop `gorm:"many2many:workgroup_monitoringshops;"`
}

func (workGroup WorkGroup) Serializer() common.H {
	var regions []common.H
	for _, region := range workGroup.Regions {
		regions = append(regions, region.Serializer())
	}
	return common.H{
		"id":      workGroup.ID,
		"name":    workGroup.Name,
		"address": workGroup.Address,
		"regions": regions,
	}
}
