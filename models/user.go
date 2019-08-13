package models

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

type User struct {
	gorm.Model
	FirstName   string       `gorm:"size:255;not null" form:"required;label:Имя"`
	LastName    string       `gorm:"size:255;" form:"required;label:Фамилия"`
	UserName    string       `gorm:"size:255;unique_index;not null" form:"required;label:Login"`
	Password    string       `gorm:"size:60;not null" form:"required;type:password;label:Пароль"`
	Email       string       `gorm:"type:varchar(100);unique_index;not null" form:"required"`
	Phone       string       `gorm:"type:varchar(17);" form:"label:Телефон"`
	Roles       []Role       `gorm:"many2many:users_roles;" form:"label:Роли для администрирования"`
	Monitorings []Monitoring `gorm:"many2many:monitorings_users;" form:"label:Мониторинги(По группам мониторинга);group_by:MonitoringGroups"`
	Online      bool         `gorm:"default:false" form:"disabled"`
	Active      bool         `gorm:"default:true" form:"type:switch;label:Активировать"`
	IsSuperUser bool         `gorm:"default:false"`
	IsStaff     bool         `gorm:"default:false" form:"label:Сотрудник"`
	TokenId     uint
	Token       Token
}

func (user User) GetFullName() string {
	return fmt.Sprintf("%s %s", user.LastName, user.FirstName)
}

func (user User) GetMonitorings() string {
	names := make([]string, 0)
	for _, m := range user.Monitorings {
		names = append(names, m.Name)
	}
	return strings.Join(names, ", ")
}

func (user *User) SetPhone(phone string) error {
	phone = strings.Trim(phone, "")
	if phone != "" {
		if matched, _ := regexp.MatchString("^\\+7\\(\\d{3}\\)-\\d{3}-\\d{2}-\\d{2}", phone); !matched {
			return errors.New("Неверно заполнен номер телефона")
		}
	}
	user.Phone = phone
	return nil
}

func (user *User) SetUserName(username string) error {
	username = strings.Trim(username, "")
	if len(username) < 3 {
		return errors.New("Login должен состоять минимум из 3х символов")
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
		return errors.New("Неверно заполнен email")
	}
	user.Email = email
	return nil
}

func (user *User) SetPassword(password string) error {
	if password = strings.Trim(password, ""); len(password) < 3 {
		return errors.New("Пароль должен состоять минимум из 3х символов")
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
		Name:   "Пользователь",
		Plural: "Пользователи",
	}
}

func (user User) Admin() types.AdminMeta {
	return types.AdminMeta{
		Preload:       []string{"Monitorings"},
		ExcludeFields: []string{"TokenId", "IsSuperUser", "Token"},
		OrderBy:       []string{"LastName", "FirstName"},
		SearchFields:  []string{"LastName", "FirstName", "Email", "Phone"},
		SortFields: []types.AdminMetaField{
			{Name: "LastName"},
			{Name: "FirstName"},
			{Name: "Email"},
		},
		ExtraFields: []types.AdminMetaField{
			{
				Name:  "GetMonitorings",
				Label: "Рабочая группа",
			},
		},
		FilterFields: []types.AdminMetaField{
			{
				Name:  "Monitorings.ID",
				Label: "По мониторингам",
			},
		},
	}
}

func (user User) String() string {
	return fmt.Sprintf("%s %s", user.LastName, user.FirstName)
}

func (user *User) CRUD(db *gorm.DB) types.CRUDManager {
	return user.Manager(db)
}

func (user *User) Manager(db *gorm.DB) *UserManager {
	return &UserManager{DB: db, self: user}
}

func (user User) Serializer() types.H {
	var roles []types.H
	var monitorings []types.H
	for _, role := range user.Roles {
		roles = append(roles, role.Serializer())
	}
	for _, m := range user.Monitorings {
		monitorings = append(monitorings, m.Serializer())
	}
	return types.H{
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"username":      user.UserName,
		"token":         user.Token.Key,
		"email":         user.Email,
		"online":        user.Online,
		"roles":         roles,
		"phone":         user.Phone,
		"active":        user.Active,
		"monitorings":   monitorings,
		"is_super_user": user.IsSuperUser,
		"is_staff":      user.IsStaff,
	}
}
