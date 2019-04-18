package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type PermissionAccess int

const (
	FORBIDDEN PermissionAccess = 0
	READ      PermissionAccess = 3
	WRITE     PermissionAccess = 5
	ACCESS    PermissionAccess = 7
)

var ChoiceAccess = map[PermissionAccess]string{
	FORBIDDEN: "Не разрешено",
	READ:      "Только для чтения",
	WRITE:     "Доступ на запись",
	ACCESS:    "Полный доступ (Чтение, Запись, Удаление)",
}

type Permission struct {
	gorm.Model
	ViewID uint
	View   Views
	Access PermissionAccess `gorm:"default:3"`
}

func (permission Permission) GetPermissionName() string {
	return ChoiceAccess[permission.Access]
}

func (permission Permission) Serializer() types.H {
	return types.H{
		"access_code": permission.Access,
		"access_name": permission.GetPermissionName(),
		"views":       permission.View.Serializer(),
	}
}
