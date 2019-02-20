package models

import "github.com/jinzhu/gorm"

type Access int

const (
	NO_ACCESS 	Access = 0
	READ 		Access = 2
	WRITE 		Access = 5
	ACCESS 		Access = 7
)

type Permission struct {
	gorm.Model
	View 	string
	Access 	Access `gorm:"default:2"`
}