package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Ware struct {
	gorm.Model
	Name        string  `gorm:"size:255" form:"label:Название товара;required"`
	Code        string  `gorm:"size:255" form:"label: Локальный код товара;required"`
	Barcode     string  `gorm:"size:255" form:"label: ШК товара;"`
	Description string  `form:"label:Описание"`
	Segment     Segment `form:"label:Сегмент"`
	SegmentId   uint
	Monitorings []Monitoring `gorm:"many2many:monitorings_wares" form:"label:Мониторинги(по типам);group_by:MonitoringType"`
}

func (ware Ware) GetSegmentName() string {
	return fmt.Sprintf("%s %s", ware.Segment.Code, ware.Segment.Name)
}

func (ware Ware) Serializer() types.H {
	return types.H{
		"id":          ware.ID,
		"name":        ware.Name,
		"code":        ware.Code,
		"description": ware.Description,
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
		SearchFields:  []string{"Name", "Barcode", "Code"},
		Preload:       []string{"Segment"},
		OrderBy:       []string{"SegmentId", "Name"},
		SortFields: []types.AdminMetaField{
			{Name: "Code"},
			{Name: "Name"},
			{Name: "Barcode"},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "GetSegmentName",
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
