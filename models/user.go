package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName	string 	`gorm:"size:255;not null" json:"first_name"`
	LastName	string 	`gorm:"size:255;not null" json:"last_name"`
	UserName	string 	`gorm:"size:255;unique_index;not null" json:"username"`
	Password    string 	`gorm:"size:60;not null" json:"password"`
	Email		string 	`gorm:"type:varchar(100);unique_index;not null" json:"email"`
	Phone		string	`gorm:"type:varchar(17);unique_index;" json:"phone"`
	Roles		[]Role 	`gorm:"many2many:user_role;" json:"role"`
	Token 		Token  	`gorm:"foreignkey:UserRefer"`
}

func (user *User) GetById(db *gorm.DB, id int64)  {
	db.Preload("Token").Preload("Role").Where("id = ?", id).First(&user)
}

func (user *User) Create(db *gorm.DB) error {
	_user := User{}

	db.Where("user_name = ?", user.UserName).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("имя пользователя %s уже занято", user.UserName)
	}

	db.Where("email = ?", user.Email).First(&_user)
	if _user.ID != 0 {
		return fmt.Errorf("email %s уже занят", user.Email)
	}

	user.Password = common.HashAndSalt(user.Password)
	db.Create(&user)
	db.NewRecord(user)

	token := Token{}
	token.Create(db, user)
	user.Token = token
	return nil
}

func (user User) Serializer() common.H {
	var roles []common.H
	for _, role := range user.Roles {
		roles = append(roles, role.Serializer())
	}
	return common.H{
		"id":			user.ID,
		"first_name": 	user.FirstName,
		"last_name": 	user.LastName,
		"username": 	user.UserName,
		"token": 		user.Token.Key,
		"email": 		user.Email,
		"roles": 		roles,
		"phone":		user.Phone,
	}
}