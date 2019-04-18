package models

import (
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

type User struct {
	gorm.Model
	FirstName   string      `gorm:"size:255;not null"`
	LastName    string      `gorm:"size:255;"`
	UserName    string      `gorm:"size:255;unique_index;not null"`
	Password    string      `gorm:"size:60;not null" json:"password"`
	Email       string      `gorm:"type:varchar(100);unique_index;not null"`
	Phone       string      `gorm:"type:varchar(17);"`
	Roles       []Role      `gorm:"many2many:users_roles;"`
	WorkGroup   []WorkGroup `gorm:"many2many:users_work_groups;"`
	Active      bool        `gorm:"default:true"`
	IsSuperUser bool        `gorm:"default:false"`
	IsStaff     bool        `gorm:"default:false"`
	TokenID     uint
	Token       Token
}

func (user *User) SetPhone(phone string) error {
	phone = strings.Trim(phone, "")
	if matched, _ := regexp.MatchString("^\\+7\\(\\d{3}\\)-\\d{3}-\\d{2}-\\d{2}", phone); !matched {
		return errors.New("not valid phone")
	}
	user.Phone = phone
	return nil
}

func (user *User) SetUserName(username string) error {
	username = strings.ToLower(strings.Trim(username, ""))
	if len(username) < 3 {
		return errors.New("not valid username")
	}
	user.UserName = username
	return nil
}

func (user *User) SetEmail(email string) error {
	email = strings.ToLower(strings.Trim(email, ""))
	matched, _ := regexp.MatchString(
		"^[a-z0-9_][a-z0-9.\\-_]{1,100}@[a-z0-9\\-_]{1,100}\\.[a-z0-9\\-_]{1,50}[a-z0-9_]$",
		email,
	)
	if !matched {
		return errors.New("not valid email")
	}
	user.Email = email
	return nil
}

func (user *User) SetPassword(password string) error {
	if password = strings.Trim(password, ""); len(password) < 3 {
		return errors.New("not valid password")
	}
	user.Password = password
	user.HashPassword()
	return nil
}

func (user *User) HashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	user.Password = string(hash)
}

func (user User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user User) Meta() types.ModelsMeta {
	return types.ModelsMeta{
		Name: "Пользователь",
	}
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
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"username":      user.UserName,
		"token":         user.Token.Key,
		"email":         user.Email,
		"roles":         roles,
		"phone":         user.Phone,
		"active":        user.Active,
		"work_groups":   workGroups,
		"is_super_user": user.IsSuperUser,
		"is_staff":      user.IsStaff,
	}
}
