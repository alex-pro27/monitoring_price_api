package models

import (
	"errors"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
)

type CompleteWareManager struct {
	*gorm.DB
	self *CompletedWare
}

func (manager *CompleteWareManager) Create(fields types.H) (err error) {
	return errors.New("No implementation")
}

func (manager *CompleteWareManager) Update(fields types.H) (err error) {
	return errors.New("No implementation")
}

func (manager *CompleteWareManager) Delete(fields types.H) (err error) {
	return errors.New("No implementation")
}
