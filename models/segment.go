package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

/**
Товарный сегмент(Сыры, колбасы, молоко)
*/
type Segment struct {
	gorm.Model
	Name  string `gorm:"size:255"`
	Code  string `gorm:"size:255"`
	Wares []Ware `gorm:"foreignkey:SegmentID"`
}

func (segment Segment) Serializer() common.H {
	var wares []common.H
	for _, ware := range segment.Wares {
		wares = append(wares, ware.Serializer())
	}
	return common.H{
		"id":    segment.ID,
		"name":  segment.Name,
		"code":  segment.Code,
		"wares": wares,
	}
}
