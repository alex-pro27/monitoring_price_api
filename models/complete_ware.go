package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"time"
)

/**
Промониторенный товар
*/
type CompletedWare struct {
	gorm.Model
	User             User
	UserID           uint
	DateUpload       time.Time
	Missing          bool
	Discount         bool
	Price            float64
	MinPrice         float64
	MaxPrice         float64
	Description      string
	Comment          string
	Ware             Ware `gorm:"auto_preload"`
	WareID           uint
	MonitoringShop   MonitoringShop
	MonitoringShopID uint
	MonitoringType   MonitoringType
	MonitoringTypeID uint
	Region           Regions
	RegionID         uint
	Barcode          string `gorm:"size:255;"`
}

func (CompletedWare) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Промониторинный товар",
		Plural: "Промониторинные товары",
	}
}

func (completeWare CompletedWare) String() string {
	return fmt.Sprintf("%s %s %s", completeWare.Ware.Code, completeWare.Ware.Name)
}
