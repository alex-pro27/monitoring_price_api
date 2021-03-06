package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Магазин который мониторят
*/
type MonitoringShop struct {
	gorm.Model
	/**
	Название
	*/
	Name string `gorm:"size:255" form:"label:Имя;required"`
	/**
	Адрес магазина
	*/
	Address string `gorm:"size:255" form:"label:Адрес"`
	/**
	Код магазина(для 1с)
	*/
	Code string `gorm:"size:255" form:"label:Код"`
	/**
	Обязательность фотографирования
	*/
	IsMustPhoto bool `form:"label:Обязательность фотографирования"`
	/**
	Мониториги
	*/
	WorkGroups []WorkGroup `gorm:"many2many:work_groups_monitoring_shops" form:"label:Рабочие группы"`
	/**
	Доступные сегменты
	*/
	Segments []Segment `gorm:"many2many:monitoring_shops_segments" form:"label:Доступные сегменты"`

	Active bool `gorm:"default:true" form:"label:Активный;type:switch"`
}

func (monitoringShop MonitoringShop) Serializer() types.H {

	return types.H{
		"id":            monitoringShop.ID,
		"code":          monitoringShop.Code,
		"name":          monitoringShop.Name,
		"address":       monitoringShop.Address,
		"is_must_photo": monitoringShop.IsMustPhoto,
	}
}

func (MonitoringShop) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Магазин для мониторинга",
		Plural: "Магазины для мониторинга",
	}
}

func (MonitoringShop) Admin() types.AdminMeta {
	return types.AdminMeta{
		SearchFields: []string{"Name", "Address"},
		OrderBy:      []string{"Code", "Name"},
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "Address"},
		},
	}
}

func (monitoringShop MonitoringShop) String() string {
	return fmt.Sprintf("%s %s %s", monitoringShop.Code, monitoringShop.Name, monitoringShop.Address)
}
