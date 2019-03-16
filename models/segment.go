package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

/**
Товарный сегмент
*/
type Segment struct {
	gorm.Model
	Name   string `gorm:"size:255"`
	Code   string `gorm:"size:255"`
	Wares  []Ware `gorm:"foreignkey:SegmentID"`
	Active bool   `gorm:"default:true"`
}

func (segment Segment) Serializer() types.H {
	var wares []types.H
	for _, ware := range segment.Wares {
		wares = append(wares, ware.Serializer())
	}
	return types.H{
		"id":     segment.ID,
		"active": segment.Active,
		"name":   segment.Name,
		"code":   segment.Code,
		"wares":  wares,
	}
}
