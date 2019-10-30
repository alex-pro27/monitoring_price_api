package models

import (
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type ViewsManager struct {
	*gorm.DB
	self *Views
}

func (manager *ViewsManager) Create(fields types.H) (err error) {
	viewType := fields["view_type"]
	if viewType != nil {
		manager.self.ViewType = ViewType(viewType.(float64))
	}
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.DB.Create(manager.self)
	manager.NewRecord(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *ViewsManager) Update(fields types.H) (err error) {
	viewType := fields["view_type"]
	if viewType != nil {
		manager.self.ViewType = ViewType(viewType.(float64))
	}
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	manager.Save(manager.self)
	helpers.SetManyToMany(manager.DB, manager.self, fields)
	return nil
}

func (manager *ViewsManager) Delete() (err error) {
	if res := manager.DB.Delete(manager.self); res.Error != nil {
		return res.Error
	}
	return nil
}
