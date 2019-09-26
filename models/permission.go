package models

import (
	"fmt"
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

var PermissionAccessChoices = map[PermissionAccess]string{
	FORBIDDEN: "Не разрешено",
	READ:      "Только для чтения",
	WRITE:     "Доступ на запись",
	ACCESS:    "Полный доступ (Чтение, Запись, Удаление)",
}

type Permission struct {
	gorm.Model
	ViewId uint
	View   Views            `gorm:"foreignkey:ViewId" form:"label:Представление"`
	Access PermissionAccess `gorm:"default:3" form:"choice:GetChoiceAccess; label:Доступ"`
}

func (Permission) GetChoiceAccess() map[PermissionAccess]string {
	return PermissionAccessChoices
}

func (permission Permission) GetPermissionName() string {
	return PermissionAccessChoices[permission.Access]
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

func (Permission) Admin() types.AdminMeta {
	return types.AdminMeta{
		ExcludeFields: []string{"ViewId"},
		Preload:       []string{"View"},
		OrderBy:       []string{"ViewId"},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "View.Name",
				Label: "Представление",
			},
			{
				Name:  "GetPermissionName",
				Label: "Доступ",
			},
		},
	}
}

func (permission *Permission) CRUD(db *gorm.DB) types.CRUDManager {
	return &PermissionManager{db, permission}
}

func (permission Permission) String() string {
	return fmt.Sprintf("%s - %s", permission.View.Name, permission.GetPermissionName())
}
