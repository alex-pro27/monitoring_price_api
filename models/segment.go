package models

import "github.com/jinzhu/gorm"

type Segment struct {
	gorm.Model
	Name 	string 		`gorm:"size:255"`
	Code 	int
	Regions []Regions 	`gorm:"many2many:segments_regions"`
	Wares 	[]Ware 		`gorm:"foreignkey:SegmentRef"`
}
