package models

import "github.com/jinzhu/gorm"

type Photos struct {
	gorm.Model
	Path            string
	CompletedWare   CompletedWare
	CompletedWareID uint
}
