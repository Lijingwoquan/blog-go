package jwt

import (
	"blog/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var mySecret []byte

func init() {
	// 从配置中读取 secret，确保在使用前已正确设置
	mySecret = []byte(viper.GetString("auth.jwt_secret"))
}

// GenToken 生成JWT
func GenToken(u *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour)

	claims := &MyClaims{
		UserID: u.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "liuzihao-blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
