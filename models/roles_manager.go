package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type RolesManager struct {
	*gorm.DB
	self *Role
}

func (manager *RolesManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	role_type := fields["role_type"]
	if role_type != nil {
		manager.self.RoleType = RoleType(role_type.(float64))
	}
	manager.DB.Create(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *RolesManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	role_type := fields["role_type"]
	if role_type != nil {
		manager.self.RoleType = RoleType(role_type.(float64))
	}
	manager.Save(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *RolesManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}
