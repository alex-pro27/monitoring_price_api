package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/common"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	FirstName string      `gorm:"size:255;not null"`
	LastName  string      `gorm:"size:255;not null"`
	UserName  string      `gorm:"size:255;unique_index;not null"`
	Password  string      `gorm:"size:60;not null" json:"password"`
	Email     string      `gorm:"type:varchar(100);unique_index;not null"`
	Phone     string      `gorm:"type:varchar(17);"`
	Roles     []Role      `gorm:"many2many:user_role;"`
	WorkGroup []WorkGroup `gorm:"many2many:user_workgroup;"`
	Active    bool        `gorm:"default:true"`
	TokenID   uint
	Token     Token
}

func (user *User) GetById(db *gorm.DB, id int) {
	db.Preload(
		"Token",
	).Preload(
		"WorkGroup",
	).Preload(
		"Roles",
	).First(
		&user, "active = true AND id = ?", id,
	)
}

func (user *User) GetByUserName(db *gorm.DB, username string) {
	db.Preload(
		"Token",
	).Preload(
		"WorkGroup",
	).Preload(
		"Roles",
	).First(
		&user, "active = true AND user_name = ?", username,
	)
}

func (user *User) GetUserByToken(db *gorm.DB, token string) {
	db.First(&user.Token, "key = ?", token)
	if user.Token.ID != 0 {
		db.Preload(
			"Roles",
		).Preload(
			"WorkGroup",
		).First(
			&user, "active = true AND token_id = ?", user.Token.ID,
		)
	}
}

func (user User) CheckPassword(password string) bool {
	return common.CompareHashAndPassword(user.Password, password)
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
	user.NewToken(db)
	user.HashPassword()
	db.Create(&user)
	db.NewRecord(user)
	return nil
}

func (user *User) HashPassword() {
	user.Password = common.HashAndSalt(user.Password)
}

func (user *User) NewToken(db *gorm.DB) {
	token := Token{}
	token.Create(db, user)
	user.TokenID = token.ID
	user.Token = token
}

func (user User) Serializer() common.H {
	var roles []common.H
	var workGroups []common.H
	for _, role := range user.Roles {
		roles = append(roles, role.Serializer())
	}
	for _, wg := range user.WorkGroup {
		workGroups = append(workGroups, wg.Serializer())
	}
	return common.H{
		"id":          user.ID,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"username":    user.UserName,
		"token":       user.Token.Key,
		"email":       user.Email,
		"roles":       roles,
		"phone":       user.Phone,
		"active":      user.Active,
		"work_groups": workGroups,
	}
}
