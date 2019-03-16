package models

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Key string `gorm:"size:32;unique_index;not null" json:"token_key"`
}
