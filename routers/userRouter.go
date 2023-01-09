package routers

import (
	"encoding/json"
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/models"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Input data error - "+err.Error(), 400)
		return
	}
	if len(u.Email) == 0 {
		http.Error(w, "email is required", 400)
		return
	}
	if len(u.Password) < 6 {
		http.Error(w, "password length must be grater than 6", 400)
		return
	}

	_, userFound, _ := db.SearchUserByEmail(u.Email)
	if userFound == true {
		http.Error(w, "User already exists", 400)
		return
	}
	_, status, err := db.SaveUser(u)
	if err != nil {
		http.Error(w, "Error on saving user "+err.Error(), 400)
		return
	}
	if status == false {
		http.Error(w, "Status was false on saving user ", 400)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
