package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Monitoring struct {
	gorm.Model
	Name             string `gorm:"size:255;not null" form:"label:Название;required"`
	MonitoringTypeId uint
	MonitoringType   MonitoringType `form:"label:Тип мониторинга"`
	Region           []Region       `gorm:"many2many:regions_monitorings;" form:"label: Регион"`
	Wares            []Ware         `gorm:"many2many:monitorings_wares" form:"label:Товары;group_by:Segment"`
	WorkGroups       []WorkGroup    `gorm:"many2many:monitorings_work_groups" form:"label:Рабочие группы"`
	Active           bool           `gorm:"default:true" form:"label:Активная;type:switch"`
}

func (monitoring Monitoring) Serializer() types.H {
	var regions []types.H
	for _, region := range monitoring.Region {
		regions = append(regions, region.Serializer())
	}
	return types.H{
		"id":      monitoring.ID,
		"name":    monitoring.Name,
		"regions": regions,
		"active":  monitoring.Active,
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
		OrderBy: []string{
			"Name",
		},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "Active"},
		},
	}
}

func (monitoring Monitoring) String() string {
	return monitoring.Name
}
