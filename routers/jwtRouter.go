package routers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/el27egs/twittor-tilotta/db"
	"github.com/el27egs/twittor-tilotta/models"
	"strings"
)

var Email string
var IDUser string

func ValidateJWT(authHeader string) (*models.Claim, bool, string, error) {
	claims := &models.Claim{}
	tokenArray := strings.Split(authHeader, "Bearer ")
	if len(tokenArray) != 2 {
		return &models.Claim{}, false, "", errors.New("jwt no encontrado")
	}
	tokenStr := strings.TrimSpace(tokenArray[1])

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("pa$$w0rd"), nil
	})
	if err != nil {
		return &models.Claim{}, false, "", err
	}
	if !tkn.Valid {
		return &models.Claim{}, false, "", errors.New("token JWT invaludi")
	}
	_, userFound, _ := db.SearchUserByEmail(claims.Email)
	if userFound {
		Email = claims.Email
		IDUser = claims.ID.Hex()
	}
	return claims, true, IDUser, nil
}
