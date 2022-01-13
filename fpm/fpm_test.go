package fpm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/team4yf/fpm-go-pkg/utils"
)

func TestToken(t *testing.T) {
	utils.InitJWTUtil("abc")
	token := GenerateToken("foo", "biz", 7200)

	checked, _ := utils.CheckToken(token.AccessToken)
	assert.True(t, checked)
}
