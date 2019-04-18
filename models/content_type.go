package models

import "github.com/jinzhu/gorm"

type ContentType struct {
	gorm.Model
	Table string  `gorm:"size:255;"`
	Views []Views `gorm:"foreignkey:ContentTypeID"`
}
