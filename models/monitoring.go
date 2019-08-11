package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Monitoring struct {
	gorm.Model
	Name             string `gorm:"size:255" form:"label:Название;required"`
	MonitoringTypeId uint
	MonitoringType   MonitoringType   `form:"label:Тип мониторинга"`
	Wares            []Ware           `gorm:"many2many:monitorings_wares" form:"label:Товары;group_by:Segment"`
	MonitoringShops  []MonitoringShop `gorm:"many2many:monitoring_shops_monitorings" form:"label:Магазины для мониторинга"`
	WorkGroups       []WorkGroup      `gorm:"many2many:work_groups_monitorings;" form:"label:Рабочие группы"`
	Active           bool             `gorm:"default:true" form:"label:Активный;type:switch"`
}

func (Monitoring) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Мониторинг",
		Plural: "Мониторинги",
	}
}

func (Monitoring) Admin() types.AdminMeta {
	return types.AdminMeta{
		SortFields: []types.AdminMetaField{
			{Name: "ID", Label: "Код"},
			{Name: "Name"},
			{Name: "UpdatedAt", ToHTML: "datetime", Label: "Дата обновления"},
			{Name: "Active"},
		},
		OrderBy: []string{"-UpdatedAt"},
	}
}

func (monitoring Monitoring) String() string {
	return fmt.Sprintf("%d %s", monitoring.ID, monitoring.Name)
}
