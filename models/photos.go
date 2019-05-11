package models

import (
	"fmt"
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
		ShortToHtml:   "image",
		ExcludeFields: []string{"CompletedWareId"},
		Preload:       []string{"CompletedWare.Ware"},
		OrderBy:       []string{"-CreatedAt"},
		ExtraFields: []types.AdminMetaField{
			{
				Label:  "Фото",
				Name:   "GetPhotoUrl",
				ToHTML: "image",
			},
			{
				Name:  "CompletedWare.Ware.Name",
				Label: "Промониторинный товар",
			},
			{
				Name:   "Model.CreatedAt",
				Label:  "Дата создания",
				ToHTML: "datetime",
			},
		},
	}
}

func (photo Photos) GetPhotoUrl() string {
	return fmt.Sprintf("/api/admin/media/%s", photo.Path)
}

func (photo Photos) String() string {
	return photo.GetPhotoUrl()
}
