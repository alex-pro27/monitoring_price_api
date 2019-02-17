package models

import "github.com/jinzhu/gorm"

type WorkGroup struct {
	gorm.Model
	Name 	string 		`gorm:"size:255"`
	Address string 		`gorm:"size:255"`
	Regions []Regions 	`gorm:"many2many:workgroup_regions;"`
	Rivals 	[]Rival  	`gorm:"many2many:workgroup_rivals;"`
}
