package models

import (
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/jinzhu/gorm"
)

type TokenManager struct {
	*gorm.DB
}

func (objects *TokenManager) NewToken(user *User) {
	token := Token{}
	if user.TokenID != 0 {
		objects.Delete(&Token{}, user.TokenID)
	}
	key := utils.GenerateToken()
	t := Token{}
	objects.First(&t, "key = ?", key)
	if t.ID != 0 {
		objects.NewToken(user)
	} else {
		token.Key = key
		objects.Create(&token)
		objects.NewRecord(token)
		user.TokenID = token.ID
		user.Token = token
	}
}
