package db

import (
	"github.com/el27egs/twittor-tilotta/models"
	"golang.org/x/crypto/bcrypt"
)

func LoginDB(email string, password string) (models.User, bool) {
	user, userFound, _ := SearchUserByEmail(email)
	if userFound == false {
		return user, false
	}
	passInput := []byte(password)
	passDB := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(passDB, passInput)
	if err != nil {
		return user, false
	}
	return user, true
}
