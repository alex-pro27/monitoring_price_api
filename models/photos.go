package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Photos struct {
	gorm.Model
	Path            string
	CompletedWareId uint `gorm:"default:null"`
	CompletedWare   CompletedWare
}

func (Photos) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Фотография",
		Plural: "Фотографии",
	}
}

func (Photos) Admin() types.AdminMeta {
	return types.AdminMeta{
		ExcludeFields: []string{"CompletedWareId"},
		Preload:       []string{"CompletedWare.Ware"},
		SortFields:    []string{"Path"},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "CompletedWare.Ware.Name",
				Label: "Промониторинный товар",
				Type:  "img",
			},
			{
				Name:  "Model.CreatedAt",
				Label: "Дата создания",
				Type:  "date",
			},
		},
	}
}

func (photo Photos) String() string {
	return photo.Path
}
