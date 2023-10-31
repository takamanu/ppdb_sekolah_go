package middleware

import (
	"ppdb_sekolah_go/constans"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userId int, name string, role int) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = userId
	claims["name"] = name
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constans.SECRET_JWT))

}
