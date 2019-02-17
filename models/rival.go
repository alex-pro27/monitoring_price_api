package models

import "github.com/jinzhu/gorm"

type Rival struct {
	gorm.Model
	Name 		string `gorm:"size:255"`
	Address 	string `gorm:"size:255"`
	code 		int
	IsMustPhoto bool
	WorkGroups  []WorkGroup `gorm:"many2many:workgroup_rivals;"`
}