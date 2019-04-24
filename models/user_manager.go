package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

type UserManager struct {
	*gorm.DB
	self *User
}

func (manager *UserManager) Create(fields types.H) (err error) {
	_user := User{}
	helpers.SetFieldsOnModel(manager.self, fields)
	manager.Where("user_name = ?", manager.self.UserName).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("имя пользователя %s уже занято", manager.self.UserName)
	}

	manager.Where("email = ?", manager.self.Email).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("email %s уже занят", manager.self.Email)
	}
	token := Token{}
	token.Manager(manager.DB).NewToken(manager.self)
	manager.DB.Create(manager.self)
	manager.NewRecord(manager.self)
	return nil
}

func (manager *UserManager) Update(fields types.H) (err error) {
	errs := helpers.SetFieldsOnModel(manager.self, fields)

	if errs != "" {
		return errors.New(errs)
	}
	_user := User{}
	manager.First(&_user, "email = ?", manager.self.Email)
	if _user.ID > 0 && _user.ID != manager.self.ID {
		return fmt.Errorf("email %s занят", manager.self.Email)
	}

	manager.self.Active = fields["active"].(bool)

	token := Token{}
	token.Manager(manager.DB).NewToken(manager.self)
	manager.Save(manager.self)
	return nil
}

func (manager *UserManager) Delete(fields types.H) (err error) {
	now := time.Now()
	manager.self.DeletedAt = &now
	manager.self.Active = false
	manager.Save(manager.self)
	return nil
}

func (manager *UserManager) GetById(id uint) *User {
	manager.Preload(
		"Token",
	).Preload(
		"WorkGroup",
	).Preload(
		"WorkGroup.Regions",
	).Preload(
		"Roles",
	).First(
		manager.self, id,
	)
	return manager.self
}

func (manager *UserManager) GetByUserName(username string) *User {
	manager.Preload(
		"Token",
	).Preload(
		"WorkGroup.Regions",
	).Preload(
		"Roles.Permissions.View",
	).First(
		manager.self, "active = true AND user_name = ? OR email = ?", username, username,
	)
	return manager.self
}

func (manager *UserManager) GetUserByToken(token string) *User {
	manager.First(&manager.self.Token, "key = ?", token)
	manager.Preload(
		"WorkGroup",
	).Preload(
		"WorkGroup.Regions",
	).Find(
		manager.self, "token_id = ?", manager.self.Token.ID,
	)
	return manager.self
}
