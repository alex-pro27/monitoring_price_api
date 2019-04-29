package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type ContentType struct {
	gorm.Model
	Table string  `gorm:"size:255;"`
	Views []Views `gorm:"foreignkey:ContentTypeID"`
}

func (contentType ContentType) String() string {
	return fmt.Sprintf(contentType.Table)
}
