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
	Name string `gorm:"size:255"`
	/**
	Адрес магазина
	*/
	Address string `gorm:"size:255"`
	/**
	Код магазина(для 1с)
	*/
	Code string `gorm:"size:255"`
	/**
	Обязательность фотографирования
	*/
	IsMustPhoto bool

	/**
	Группа мониторинга
	*/
	WorkGroup []WorkGroup `gorm:"many2many:work_groups_monitoring_shops;"`

	/**
	Сегменты
	*/
	Segments []Segment `gorm:"many2many:monitoring_shops_segments;"`
}

func (monitoringShop MonitoringShop) Serializer() types.H {
	var segmentIDX []uint
	for _, segment := range monitoringShop.Segments {
		segmentIDX = append(segmentIDX, segment.ID)
	}
	return types.H{
		"id":            monitoringShop.ID,
		"code":          monitoringShop.Code,
		"name":          monitoringShop.Name,
		"address":       monitoringShop.Address,
		"is_must_photo": monitoringShop.IsMustPhoto,
		"segments":      segmentIDX,
	}
}

func (MonitoringShop) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Магазин для мониторинга",
		Plural: "Магазины для мониторинга",
	}
}

func (monitoringShop MonitoringShop) String() string {
	return fmt.Sprintf("%s %s %s", monitoringShop.Code, monitoringShop.Name, monitoringShop.Address)
}
