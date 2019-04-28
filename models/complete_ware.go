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
	UserId           uint
	DateUpload       time.Time
	Missing          bool
	Discount         bool
	Price            float64
	MinPrice         float64
	MaxPrice         float64
	Description      string
	Comment          string
	Ware             Ware
	WareId           uint
	MonitoringShop   MonitoringShop
	MonitoringShopId uint
	MonitoringType   MonitoringType
	MonitoringTypeId uint
	Region           Regions
	RegionId         uint   `gorm:"default:null"`
	Barcode          string `gorm:"size:255;"`
	Photos           []Photos
}

func (CompletedWare) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Промониторинный товар",
		Plural: "Промониторинные товары",
	}
}

func (completeWare CompletedWare) String() string {
	return fmt.Sprintf("%s %s", completeWare.Ware.Code, completeWare.Ware.Name)
}

func (completeWare *CompletedWare) Manager(db *gorm.DB) *CompleteWareManager {
	return &CompleteWareManager{db, completeWare}
}

func (CompletedWare) Admin() types.AdminMeta {
	return types.AdminMeta{
		Preload: []string{"Ware"},
	}
}
