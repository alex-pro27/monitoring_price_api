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

type Permission struct {
	gorm.Model
	ViewID uint
	View   Views
	Access PermissionAccess `gorm:"default:3" option:"choice:GetChoiceAccess"`
}

func (Permission) GetChoiceAccess() map[PermissionAccess]string {
	return map[PermissionAccess]string{
		FORBIDDEN: "Не разрешено",
		READ:      "Только для чтения",
		WRITE:     "Доступ на запись",
		ACCESS:    "Полный доступ (Чтение, Запись, Удаление)",
	}

}

func (permission Permission) GetPermissionName() string {
	return permission.GetChoiceAccess()[permission.Access]
}

func (permission Permission) Serializer() types.H {
	return types.H{
		"access_code": permission.Access,
		"access_name": permission.GetPermissionName(),
		"views":       permission.View.Serializer(),
	}
}

func (Permission) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Разрешение",
		Plural: "Разрешения",
	}
}

func (permission Permission) String() string {
	return permission.GetPermissionName()
}
