package utils

import "github.com/dgrijalva/jwt-go"

// GenerateToken 生成Token
func GenerateToken(mapClaims jwt.MapClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(key))
}

//  验证token
// func checkToken(uid int64, token *jwt.Token) bool {
// 	tokens, _ := token.SignedString([]byte(JWTKey))
// 	redisToken, _ := GetMemberToken(uid)
// 	if tokens != redisToken {
// 		return false
// 	}
// 	return true
// }

//  用户登录请求取出token
//  token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
// 	return []byte(JWTKey), nil
//  })
//  if err == nil && token.Valid {
// 	tokenMap := token.Claims.(jwt.MapClaims)
// 	uidStr := tokenMap["uid"].(string)
// 	uid, _ := strconv.ParseInt(uidStr,10,64)

// 	if !checkToken(uid, token) {
// 	   // 验证token 是否合法
// 	   base.ErrorResponse(w, http.StatusUnauthorized, "Authorization Is Invalid")
// 	   return
// 	}
//  }
