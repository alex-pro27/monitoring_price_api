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
	Key string `gorm:"size:32;unique_index;not null" json:"token_key"`
}

func (token *Token) Create(db *gorm.DB, user *User) {
	key := token.generate()
	t := Token{}
	db.First(&t, "key = ?", key)
	if t.ID != 0 {
		token.Create(db, user)
	} else {
		token.Key = key
		db.Create(&token)
		db.NewRecord(token)
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
