package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var (
	key []byte
)

//InitJWTUtil init the jwt key
func InitJWTUtil(sceret string) {
	key = []byte(sceret)
}

//GenerateToken 生成Token
func GenerateToken(mapClaims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(key))
}

//CheckToken  验证token
func CheckToken(tokenString string) (bool, *jwt.Token) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false, nil
	}
	return token.Valid, token
}
