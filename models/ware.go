package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Ware struct {
	gorm.Model
	Name           string  `gorm:"size:255" form:"label:Название товара;required"`
	Code           string  `gorm:"size:255" form:"label: Локальный код товара;required"`
	Barcode        string  `gorm:"size:255" form:"label: ШК товара;required"`
	Description    string  `form:"label:Описание"`
	Segment        Segment `form:"label:Сегмент"`
	SegmentId      uint
	Active         bool             `gorm:"default:true" form:"label:Активный;type:switch"`
	MonitoringType []MonitoringType `gorm:"many2many:wares_monitoring_types"`
}

func (ware Ware) Serializer() types.H {
	return types.H{
		"id":          ware.ID,
		"name":        ware.Name,
		"code":        ware.Code,
		"description": ware.Description,
		"active":      ware.Active,
	}
}

func (Ware) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Товар",
		Plural: "Товары",
	}
}

func (Ware) Admin() types.AdminMeta {
	return types.AdminMeta{
		ExcludeFields: []string{"SegmentId"},
		SearchFields:  []string{"Name", "Barcode"},
		Preload:       []string{"Segment"},
		OrderBy:       []string{"SegmentId", "Name"},
		SortFields:    []string{"Name", "Barcode", "Active"},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "Segment.Name",
				Label: "Сегмент",
			},
		},
		FilterFields: []types.AdminMetaField{
			{
				Name:  "Segment.Name",
				Label: "Сегмент",
			},
		},
	}
}

func (ware Ware) String() string {
	return fmt.Sprintf("%s %s", ware.Code, ware.Name)
}
