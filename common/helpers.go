package common

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

func HandlerError(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

func WatchFile(isChanged chan bool, filePath string) {
	defer fmt.Println("exit watch file")
	initialStat, err := os.Stat(filePath)
	HandlerError(err)
	go func () {
		for {
			stat, err := os.Stat(filePath)
			HandlerError(err)
			if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				initialStat, err = os.Stat(filePath)
				HandlerError(err)
				isChanged <- true
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func RaiseProcess(sig os.Signal) {
	p, err := os.FindProcess(os.Getpid())
	HandlerError(err)
	HandlerError(p.Signal(sig))
}

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func CompareHashAndPassword(hashedPassword string, password string) bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func IndexOf(element interface{}, data []interface{}) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
}

func Unique(intSlice []uint) []uint {
	keys := make(map[uint]bool)
	var list []uint
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}