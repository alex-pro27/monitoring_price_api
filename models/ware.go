package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"strings"
)

type Ware struct {
	gorm.Model
	Name           		string  				`gorm:"size:255" form:"label:Название товара;required"`
	Code           		string  				`gorm:"size:255" form:"label: Локальный код товара;required"`
	Barcode        		string  				`gorm:"size:255" form:"label: ШК товара;"`
	Description    		string  				`form:"label:Описание"`
	Segment        		Segment 				`form:"label:Сегмент"`
	SegmentId      		uint
	Active         		bool             		`gorm:"default:true" form:"label:Активный;type:switch"`
	MonitoringType 		[]MonitoringType 		`gorm:"many2many:wares_monitoring_types" form:"label:Тип мониторинга"`
	MonitoringShops 	[]MonitoringShop 		`gorm:"many2many:monitoring_shops_wares" form:"label:Магазины для мониторинга"`
}

func (ware Ware) GetMonitoringType() string {
	var names []string
	for _, mt := range ware.MonitoringType {
		names = append(names, mt.Name)
	}
	return strings.Join(names, ", ")
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
		Preload:       []string{"Segment", "MonitoringType"},
		OrderBy:       []string{"SegmentId", "Name"},
		SortFields: []types.AdminMetaField{
			{Name: "UpdatedAt", ToHTML: "datetime", Label: "Дата обновления"},
			{Name: "Name"},
			{Name: "Barcode"},
			{Name: "Active"},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "GetSegmentName",
				Label: "Сегмент",
			},
			{
				Name:  "GetMonitoringType",
				Label: "Типы мониторинга",
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
