package routers

import (
	"encoding/json"
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/jwt"
	"github.com/el27egs/twittor-tilotta/models"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Usuario y/o password incorrectos "+err.Error(), 400)
		return
	}
	if len(u.Email) == 0 {
		http.Error(w, "email is required", 400)
		return
	}
	user, userFound := db.LoginDB(u.Email, u.Password)
	if userFound == false {
		http.Error(w, "Usuario y/o password incorrectos ", 400)
		return
	}
	jwt, err := jwt.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Error on genratin JWT token "+err.Error(), 400)
		return
	}
	resp := models.LoginResponse{Token: jwt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	expTime := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   jwt,
		Expires: expTime,
	})
}
