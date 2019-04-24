package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Photos struct {
	gorm.Model
	Path            string
	CompletedWare   CompletedWare
	CompletedWareID uint
}

func (Photos) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Фотография",
		Plural: "Фотографии",
	}
}

func (photo Photos) String() string {
	return fmt.Sprintf("<img alt='%s' src='%s'/>", photo.CompletedWare.Ware.Name, photo.Path)
}
