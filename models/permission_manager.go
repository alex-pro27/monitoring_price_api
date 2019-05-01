package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type PermissionManager struct {
	*gorm.DB
	self *Permission
}

func (manager *PermissionManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	access := fields["access"]
	if access != nil {
		manager.self.Access = PermissionAccess(access.(float64))
	}
	manager.DB.Create(manager.self)
	manager.NewRecord(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *PermissionManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	access := fields["access"]
	if access != nil {
		manager.self.Access = PermissionAccess(access.(float64))
	}
	manager.Save(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *PermissionManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}
