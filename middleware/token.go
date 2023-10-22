package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key") // 用于签署和验证 Token 的密钥

type Claims struct {
	UserID   int    // 用户ID
	Username string // 用户名
	jwt.StandardClaims
}

// 生成token
func GenerateToken(userID int, username string) string {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 设置 Token 有效期
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

// 验证token
func VerifyToken(tokenString string) (*Claims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // 验证 Token
	})

	// 验证 Token 是否有效
	if err != nil {
		return nil, err // Token 无效
	}

	// 检查 Token 是否过期
	if !token.Valid {
		return nil, errors.New("Token 已过期")
	}

	// Token 有效，返回用户信息
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("Token 格式不正确")
}
