package utils

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	InitJWTUtil("abc")
	tokenStr, err := GenerateToken(&jwt.MapClaims{
		"name":  "abc",
		"title": "Admin",
		"role":  "admin",
		"email": "ffff",
		"uid":   1,
	})

	assert.Nil(t, err, "should not err")
	assert.NotNil(t, tokenStr, "should not nil")
	ok, token := CheckToken(tokenStr)
	assert.True(t, ok, "should be true")
	claims := token.Claims.(jwt.MapClaims)
	assert.Equal(t, "abc", claims["name"])
	assert.Equal(t, "role", claims["admin"])
}
