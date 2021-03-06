package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type CompleteWareManager struct {
	*gorm.DB
	self *CompletedWare
}

func (manager *CompleteWareManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	if res := manager.DB.FirstOrCreate(manager.self, manager.self); res.Error != nil {
		return errors.New("Ошибка добавления записи")
	}
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *CompleteWareManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	if res := manager.Save(manager.self); res.Error != nil {
		return errors.New("Ошибка обновления записи")
	}
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *CompleteWareManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}
