package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type RoleType int

const (
	IS_USER RoleType = iota
	IS_MANAGER
	IS_ADMIN
)

var RoleTypeChoices = map[RoleType]string{
	IS_USER:    "Пользователь",
	IS_MANAGER: "Менеджер",
	IS_ADMIN:   "Админ",
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:255;not null" form:"label:Название;required"`
	RoleType    RoleType     `gorm:"default:0" form:"choice:GetChoiceRoleType;label:Тип пользователя"`
	Permissions []Permission `gorm:"many2many:roles_permissions;" form:"label:Разрешение"`
	Users       []User       `gorm:"many2many:users_roles;" form:"label:Пользователи"`
}

func (role Role) GetChoiceRoleType() map[RoleType]string {
	return RoleTypeChoices
}

func (role Role) GetRoleTypeName() string {
	return RoleTypeChoices[role.RoleType]
}

func (role *Role) CRUD(db *gorm.DB) types.CRUDManager {
	return &RolesManager{db, role}
}

func (role Role) Serializer() types.H {
	var permissions []types.H
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Serializer())
	}
	return types.H{
		"id":             role.ID,
		"name":           role.Name,
		"permissions":    permissions,
		"role_type":      role.RoleType,
		"role_type_name": role.GetRoleTypeName(),
	}
}

func (Role) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name:   "Роль",
		Plural: "Роли",
	}
}

func (role Role) String() string {
	return role.Name
}
