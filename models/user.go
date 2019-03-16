package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	FirstName string      `gorm:"size:255;not null"`
	LastName  string      `gorm:"size:255;not null"`
	UserName  string      `gorm:"size:255;unique_index;not null"`
	Password  string      `gorm:"size:60;not null" json:"password"`
	Email     string      `gorm:"type:varchar(100);unique_index;not null"`
	Phone     string      `gorm:"type:varchar(17);"`
	Roles     []Role      `gorm:"many2many:users_roles;"`
	WorkGroup []WorkGroup `gorm:"many2many:users_work_groups;"`
	Active    bool        `gorm:"default:true"`
	TokenID   uint
	Token     Token
}

func (user *User) SetPassword(password string) {
	user.Password = password
	user.HashPassword()
}

func (user *User) HashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(hash)
}

func (user User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user User) Serializer() types.H {
	var roles []types.H
	var workGroups []types.H
	for _, role := range user.Roles {
		roles = append(roles, role.Serializer())
	}
	for _, wg := range user.WorkGroup {
		workGroups = append(workGroups, wg.Serializer())
	}
	return types.H{
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
