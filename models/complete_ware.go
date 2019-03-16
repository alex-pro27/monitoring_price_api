package models

import (
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
	Ware             Ware
	WareID           uint
	MonitoringShop   MonitoringShop
	MonitoringShopID uint
	MonitoringType   MonitoringType
	MonitoringTypeID uint
	Region           Regions
	RegionID         uint
	Barcode          string `gorm:"size:255;"`
}
