package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Photos struct {
	gorm.Model
	Path            string
	CompletedWareId uint
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
	}
}

func (photo Photos) String() string {
	return fmt.Sprintf(
		"<img alt='' src='%s'/>",
		photo.Path,
	)
}
