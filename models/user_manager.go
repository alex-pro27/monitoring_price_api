package models

import (
	"encoding/json"
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/helpers"
	"github.com/alex-pro27/monitoring_price_api/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

type UserManager struct {
	*gorm.DB
	self *User
}

func (manager *UserManager) Create(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	errs := make(map[string]string)
	if res := manager.First(&User{}, "user_name = ?", manager.self.UserName); !res.RecordNotFound() {
		errs["user_name"] = fmt.Sprintf("Имя пользователя %s уже занято", manager.self.UserName)
	}

	if res := manager.First(&User{}, "email = ?", manager.self.Email); !res.RecordNotFound() {
		errs["email"] = fmt.Sprintf("Email %s уже занят", manager.self.Email)
	}

	if len(errs) > 0 {
		message, _ := json.Marshal(errs)
		return errors.New(string(message))
	}

	workGroupIDX := fields["work_group"]
	var workGroups []WorkGroup
	if reflect.ValueOf(workGroupIDX).Kind() == reflect.Slice {
		manager.Find(&workGroups, "id IN (?)", workGroupIDX)
		manager.self.WorkGroup = workGroups
	}

	rolesIDX := fields["roles"]
	var roles []Role
	if reflect.ValueOf(rolesIDX).Kind() == reflect.Slice {
		manager.Find(&roles, "id IN (?)", rolesIDX)
		manager.self.Roles = roles
	}

	(&Token{}).Manager(manager.DB).NewToken(manager.self)
	manager.DB.Create(manager.self)
	manager.NewRecord(manager.self)
	return nil
}

func (manager *UserManager) Update(fields types.H) (err error) {
	if err = helpers.SetFieldsForModel(manager.self, fields); err != nil {
		return err
	}
	errs := make(map[string]string)
	res := manager.First(&User{}, "email = ? and not id = ?", manager.self.Email, manager.self.ID)
	if !res.RecordNotFound() {
		errs["email"] = fmt.Sprintf("Email %s занят", manager.self.Email)
	}

	if len(errs) > 0 {
		message, _ := json.Marshal(errs)
		return errors.New(string(message))
	}

	workGroupIDX := fields["work_group"]
	var workGroups []WorkGroup
	if reflect.ValueOf(workGroupIDX).Kind() == reflect.Slice {
		manager.Find(&workGroups, "id IN (?)", workGroupIDX)
		manager.self.WorkGroup = workGroups
	}

	rolesIDX := fields["roles"]
	var roles []Role
	if reflect.ValueOf(rolesIDX).Kind() == reflect.Slice {
		manager.Find(&roles, "id IN (?)", rolesIDX)
		manager.self.Roles = roles
	}

	manager.self.Active = fields["active"].(bool)
	(&Token{}).Manager(manager.DB).NewToken(manager.self)
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
	).First(
		manager.self, "active = true AND user_name = ? OR email = ?", username, username,
	)
	return manager.self
}

func (manager *UserManager) GetUserByToken(token string) *User {
	manager.First(&manager.self.Token, "key = ?", token)
	manager.Preload(
		"WorkGroup.Regions",
	).Find(
		manager.self, "token_id = ?", manager.self.Token.ID,
	)
	return manager.self
}