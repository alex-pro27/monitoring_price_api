package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type PeriodManager struct {
	*gorm.DB
	self *Period
}

func (manager *PeriodManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.DB.Create(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *PeriodManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.Save(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *PeriodManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}
