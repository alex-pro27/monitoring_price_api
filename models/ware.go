package models

import "github.com/jinzhu/gorm"

type Ware struct {
	gorm.Model
	Name 		string `gorm:"size:255"`
	Code 		int
	Price 		float32
	MinPrice 	float32
	MaxPrice 	float32
	Description string
	SegmentRef 	uint
}
