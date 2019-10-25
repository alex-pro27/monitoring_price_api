package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type WorkGroup struct {
	gorm.Model
	Name            string `gorm:"size:255" form:"label:Название"`
	RegionId        uint
	Region          Region           `form:"label:Регион"`
	Monitorings     []Monitoring     `gorm:"many2many:monitorings_work_groups" form:"label:Мониторинги"`
	Users           []User           `gorm:"many2many:work_groups_users" form:"label:Пользователи"`
	MonitoringShops []MonitoringShop `gorm:"many2many:work_groups_monitoring_shops" form:"label:Магазины для мониторинга"`
}

func (workGroup WorkGroup) GetRegionName() string {
	return workGroup.Region.Name
}

func (WorkGroup) Admin() types.AdminMeta {
	return types.AdminMeta{
		ExcludeFields: []string{"RegionId"},
		SearchFields:  []string{"Name"},
		Preload:       []string{"Region"},
		OrderBy:       []string{"Name"},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "GetRegionName",
				Label: "Регион",
			},
		},
	}
}

func (WorkGroup) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Рабочая группа",
		Plural: "Группа мониторинга",
	}
}

func (workGroup WorkGroup) String() string {
	return workGroup.Name
}

func (workGroup WorkGroup) Serializer() types.H {
	return types.H{
		"id":   workGroup.ID,
		"name": workGroup.Name,
	}
}
