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
	UserId           uint
	User             User
	DateUpload       time.Time
	Missing          bool
	Discount         bool
	Price            float64
	MinPrice         float64
	MaxPrice         float64
	Description      string
	Comment          string
	WareId           uint
	Ware             Ware
	MonitoringShopId uint
	MonitoringShop   MonitoringShop
	MonitoringTypeId uint
	MonitoringType   MonitoringType
	RegionId         uint `gorm:"default:null"`
	Region           Regions
	Photos           []Photos `gorm:"foreignkey:CompletedWareId"`
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

func (comleteWare *CompletedWare) CRUD(db *gorm.DB) types.CRUDManager {
	return comleteWare.Manager(db)
}

func (completeWare *CompletedWare) Manager(db *gorm.DB) *CompleteWareManager {
	return &CompleteWareManager{db, completeWare}
}

func (CompletedWare) Admin() types.AdminMeta {
	return types.AdminMeta{
		Preload: []string{"Ware"},
	}
}
