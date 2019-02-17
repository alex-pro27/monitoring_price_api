package models

import "github.com/jinzhu/gorm"

type Regions struct {
	gorm.Model
	Name 		string		`gorm:"size:255"`
	Segments 	[]Segment 	`gorm:"many2many:segments_regions;"`
	WorkGroups  []WorkGroup `gorm:"many2many:workgroup_regions;"`
}
