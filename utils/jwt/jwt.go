package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

var mySecret = []byte("ZzAa1Ks!u!+")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

const TokenExpireDuration = time.Hour * 24

// GetToken 生成access token
func GetToken(userID int64) (aToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    viper.GetString("appName"),                           // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}

// RefreshToken 刷新AccessToken
//func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
//	// refresh token无效直接返回
//	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
//		return
//	}
//
//	// 从旧access token中解析出claims数据
//	var claims MyClaims
//	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
//	v, _ := err.(*jwt.ValidationError)
//
//	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
//	if v.Errors == jwt.ValidationErrorExpired {
//		return GetToken(claims.UserID)
//	}
//	return
//}
