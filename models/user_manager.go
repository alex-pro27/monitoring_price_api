package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserManager struct {
	*gorm.DB
}

func (objects *UserManager) Create(user *User) (*User, error) {
	_user := User{}
	objects.Where("user_name = ?", user.UserName).First(&_user)
	if _user.ID != 0 {
		return user, fmt.Errorf("имя пользователя %s уже занято", user.UserName)
	}

	objects.Where("email = ?", user.Email).First(&_user)
	if _user.ID != 0 {
		return user, fmt.Errorf("email %s уже занят", user.Email)
	}
	tokenManager := TokenManager{objects.DB}
	tokenManager.NewToken(user)
	user.HashPassword()
	objects.DB.Create(&user)
	objects.NewRecord(user)
	return user, nil
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
		&user, "active = true AND id = ?", id,
	)
	return user
}

func (objects *UserManager) GetByUserName(username string) User {
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
		&user, "active = true AND user_name = ?", username,
	)
	return user
}

func (objects *UserManager) GetUserByToken(token string) User {
	user := User{}
	objects.First(&user.Token, "key = ?", token)
	objects.Preload(
		"Roles",
	).Preload(
		"WorkGroup",
	).Preload(
		"WorkGroup.Regions",
	).Find(
		&user, "token_id = ?", user.Token.ID,
	)
	return user
}
