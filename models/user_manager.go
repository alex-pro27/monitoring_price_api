package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserManager struct {
	*gorm.DB
}

func (objects *UserManager) Create(user *User) error {
	_user := User{}
	objects.Where("user_name = ?", user.UserName).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("имя пользователя %s уже занято", user.UserName)
	}

	objects.Where("email = ?", user.Email).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("email %s уже занят", user.Email)
	}
	tokenManager := TokenManager{objects.DB}
	tokenManager.NewToken(user)
	objects.DB.Create(user)
	objects.NewRecord(user)
	return nil
}

func (objects *UserManager) GetById(id uint) User {
	user := User{}
	objects.Preload(
		"Token",
	).Preload(
		"WorkGroup",
	).Preload(
		"WorkGroup.Regions",
	).Preload(
		"Roles",
	).First(
		&user, id,
	)
	return user
}

func (objects *UserManager) GetByUserName(username string) User {
	user := User{}
	objects.Preload(
		"Token",
	).Preload(
		"WorkGroup.Regions",
	).Preload(
		"Roles.Permissions.View",
	).First(
		&user, "active = true AND user_name = ? OR email = ?", username, username,
	)
	return user
}

func (objects *UserManager) GetUserByToken(token string) User {
	user := User{}
	objects.First(&user.Token, "key = ?", token)
	objects.Preload(
		"WorkGroup",
	).Preload(
		"WorkGroup.Regions",
	).Find(
		&user, "token_id = ?", user.Token.ID,
	)
	return user
}
