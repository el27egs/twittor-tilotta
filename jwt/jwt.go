package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/el27egs/twittor-tilotta/models"
	"time"
)

func GenerateJWT(user models.User) (string, error) {

	payload := jwt.MapClaims{
		"_id":       user.ID.Hex(),
		"name":      user.Name,
		"lastName":  user.LastName,
		"birthday":  user.Birthday,
		"email":     user.Email,
		"avatar2":   user.Avatar,
		"banner":    user.Banner,
		"biography": user.Biography,
		"location":  user.Location,
		"web":       user.Web,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString([]byte("pa$$w0rd"))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
