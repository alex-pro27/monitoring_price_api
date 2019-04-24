package models

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Key string `gorm:"size:32;unique_index;not null" json:"token_key"`
}

func (token *Token) Manager(db *gorm.DB) *TokenManager {
	return &TokenManager{db, token}
}

func (token Token) String() string {
	return token.Key
}
