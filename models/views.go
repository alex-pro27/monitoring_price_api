package models

import "github.com/jinzhu/gorm"

type Views struct {
	gorm.Model
	Name       string       `gorm:"size:255"`
	Permission []Permission `gorm:"many2many:views_permissions;"`
	ParentID   uint
	Parent     *Views
	RoutePath  string  `gorm:"size:255"`
	Children   []Views `gorm:"foreignkey:ParentID"`
}
