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
	Monitorings []Monitoring `gorm:"many2many:monitorings_monitoring_shops" form:"label:Мониторинги;group_by:MonitoringGroups"`
	/**
	Товары для мониторига
	*/
	Wares []Ware `gorm:"many2many:monitoring_shops_wares" form:"label:Товары;group_by:Segment"`

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
		SortFields: []types.AdminMetaField{
			{Name: "Name"},
			{Name: "Address"},
		},
	}
}

func (monitoringShop MonitoringShop) String() string {
	return fmt.Sprintf("%s %s %s", monitoringShop.Code, monitoringShop.Name, monitoringShop.Address)
}
