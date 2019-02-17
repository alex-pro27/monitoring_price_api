package common

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)


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