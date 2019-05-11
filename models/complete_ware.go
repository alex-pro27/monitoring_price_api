package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

/**
Промониторенный товар
*/
type CompletedWare struct {
	gorm.Model
	UserId           uint           `gorm:"default:null"`
	User             User           `form:"label:Пользователь"`
	DateUpload       time.Time      `form:"label:Дата выгрузки"`
	Missing          bool           `form:"label:Отсутсвует"`
	Discount         bool           `form:"label:Скидка"`
	Price            float64        `form:"label:Цена"`
	MinPrice         float64        `form:"label:Минимальная цена"`
	MaxPrice         float64        `form:"label:Максимальная цена"`
	Comment          string         `form:"label:Комментарий"`
	WareId           uint           `gorm:"default:null"`
	Ware             Ware           `form:"label:Товар"`
	MonitoringShopId uint           `gorm:"default:null"`
	MonitoringShop   MonitoringShop `form:"label:Магазин"`
	MonitoringTypeId uint           `gorm:"default:null"`
	MonitoringType   MonitoringType `form:"label:Тип мониторинга"`
	RegionId         uint           `gorm:"default:null"`
	Region           Regions        `form:"label:Регион"`
	Photos           []Photos       `gorm:"foreignkey:CompletedWareId" form:"label:Фотографии"`
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

func (completedWare CompletedWare) GetPhotos() string {
	var photos []string
	for _, photo := range completedWare.Photos {
		photos = append(photos, photo.GetPhotoUrl())
	}
	return strings.Join(photos, ",")
}

func (CompletedWare) Admin() types.AdminMeta {
	return types.AdminMeta{
		Preload: []string{"Ware", "User", "MonitoringShop", "MonitoringType", "Region", "Photos"},
		OrderBy: []string{"-DateUpload"},
		SortFields: []types.AdminMetaField{
			{
				Name:   "DateUpload",
				ToHTML: "datetime",
			},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "Ware.Name",
				Label: "Товар",
			},
			{
				Name:  "User.GetFullName",
				Label: "Пользователь",
			},
			{
				Name:  "MonitoringShop.Name",
				Label: "Магазин",
			},
			{
				Name:  "MonitoringType.Name",
				Label: "Тип мониторинга",
			},
			{
				Name:  "Region.Name",
				Label: "Регион",
			},
			{
				Name:   "GetPhotos",
				Label:  "Фотографии",
				ToHTML: "image",
			},
		},
	}
}
