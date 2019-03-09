package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

type Ware struct {
	gorm.Model
	Name        string `gorm:"size:255"`
	Code        string `gorm:"size:255"`
	Description string
	Segment     Segment
	SegmentID   uint
	Active      bool `gorm:"default:true"`
}

func (ware Ware) Serializer() common.H {
	return common.H{
		"id":          ware.ID,
		"name":        ware.Name,
		"code":        ware.Code,
		"description": ware.Description,
		"active":      ware.Active,
	}
}
