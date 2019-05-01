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
	Wares  []Ware `gorm:"foreignkey:SegmentId"`
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

func (Segment) Admin() types.AdminMeta {
	return types.AdminMeta{
		OrderBy:    []string{"Name"},
		SortFields: []string{"Name", "Code"},
	}
}

func (Segment) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Сегмент",
		Plural: "Сегменты",
	}
}

func (segment Segment) String() string {
	return segment.Name
}
