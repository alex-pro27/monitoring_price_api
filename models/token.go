package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

type Token struct {
	gorm.Model
	Key 		string `gorm:"size:32;unique_index;not null" json:"token_key"`
	UserID 	uint
}

func (token *Token) Create(db *gorm.DB, user *User) {
	token.Key = token.generate()
	token.UserID = user.ID
	db.Create(&token)
	if db.NewRecord(token) {
		token.Create(db, user)
	}
}

func (token Token) generate() string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(strconv.FormatInt(time.Now().Unix(), 10)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hash to store:", string(hash))
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}