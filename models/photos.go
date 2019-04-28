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
	CompletedWareId uint
}

func (Photos) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Фотография",
		Plural: "Фотографии",
	}
}

func (photo Photos) String() string {
	return fmt.Sprintf("<img style='width:100px;height:100px;border-radius:50px;' alt='%s' src='%s'/>", photo.CompletedWare.Ware.Name, photo.Path)
}
