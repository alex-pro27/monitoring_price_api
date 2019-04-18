package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:255;not null"`
	Permissions []Permission `gorm:"many2many:roles_permissions;"`
}

func (role Role) Serializer() types.H {
	var permissions []types.H
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Serializer())
	}
	return types.H{
		"id":          role.ID,
		"name":        role.Name,
		"permissions": permissions,
	}
}
