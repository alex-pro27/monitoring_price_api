package models

import (
	"github.com/alex-pro27/monitoring_price_api/utils"
	"github.com/jinzhu/gorm"
)

type TokenManager struct {
	*gorm.DB
	self *Token
}

func (manager *TokenManager) NewToken(user *User) {
	if user.TokenId != 0 {
		manager.Delete(&Token{}, user.TokenId)
	}
	key := utils.GenerateToken()
	t := Token{}
	manager.First(&t, "key = ?", key)
	if t.ID != 0 {
		manager.NewToken(user)
	} else {
		manager.self.Key = key
		manager.Create(manager.self)
		manager.NewRecord(*manager.self)
		user.TokenId = manager.self.ID
		user.Token = *manager.self
	}
}
