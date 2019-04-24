package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type PeriodManager struct {
	*gorm.DB
	self *Period
}

func (manager *PeriodManager) Create(fields types.H) (err error) {
	return errors.New("No implementation")
}

func (manager *PeriodManager) Update(fields types.H) (err error) {
	return errors.New("No implementation")
}

func (manager *PeriodManager) Delete(fields types.H) (err error) {
	return errors.New("No implementation")
}
