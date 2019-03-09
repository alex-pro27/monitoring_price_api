package models

import (
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:255;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

func (role Role) Serializer() common.H {
	return common.H{
		"name": role.Name,
	}
}
