package models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.DB
	User             User
	UserId           uint
	MonitoringShop   MonitoringShop
	MonitoringShopId uint
}
