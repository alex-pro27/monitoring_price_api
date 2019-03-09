package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
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
	WorkGroup []WorkGroup `gorm:"many2many:workgroup_monitoringshops;"`

	/**
	Сегменты
	*/
	Segments []Segment `gorm:"many2many:monitoringshops_segments;"`
}

func (monitoringShop MonitoringShop) Serializer() common.H {
	var segmentIDX []uint
	for _, segment := range monitoringShop.Segments {
		segmentIDX = append(segmentIDX, segment.ID)
	}
	return common.H{
		"id":            monitoringShop.ID,
		"code":          monitoringShop.Code,
		"name":          monitoringShop.Name,
		"address":       monitoringShop.Address,
		"is_must_photo": monitoringShop.IsMustPhoto,
		"segments":      segmentIDX,
	}
}
