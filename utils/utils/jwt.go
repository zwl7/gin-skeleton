package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/dgrijalva/jwt-go"
)

//JWT 回话消息
// type JWT struct {
// 	Tag        string `json:"tag"`
// 	RequestID  string `json:"request_id"`
// 	BusinessID int64  `json:"business_id"`
// 	UserID     int64  `json:"user_id"`
// 	IsSuper    bool   `json:"is_super"`
// }

const (

	//ECDSAKeyD ES256 keys
	ECDSAKeyD = "CCFDFDC9C2572D15C639D07E3C6C8804A1E941B13F5D10C7297A2DFAA70E6393"
	//ECDSAKeyX ...
	ECDSAKeyX = "EE4C3E11EB1BF081CFD4B5CCC482E069BFBECA07D566238F29191716319B809E"
	//ECDSAKeyY ...
	ECDSAKeyY = "A40CCD993EC355326588E2A9E202C24A2D5D1BE5128B19885FD9F2C4155C3EF1"

	//SignedKey HS256 signed key
	SignedKey = "papa20140924"
)

// JwtGetEStoken 获取签名算法为ES256的token
// 该token的内容只有Redis的key,用于保存用户的登录状态
func JwtGetEStoken(redisValue string) string {
	keyD := new(big.Int)
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyD.SetString(ECDSAKeyD, 16)
	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	claims := jwt.MapClaims{
		"redisValue": redisValue,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	privateKey := ecdsa.PrivateKey{D: keyD, PublicKey: publicKey}

	ss, err := token.SignedString(&privateKey)
	if err != nil {
		fmt.Println("ES256的token生成签名错误,err:", err)
		return ""
	}
	return ss
}

// JwtGetHStoken 获取签名算法为HS256的token
func JwtGetHStoken(tokenFirst string, user map[string]string, signedKey ...string) string {
	claims := jwt.MapClaims{
		"tokenES": tokenFirst,
	}
	for _k, _v := range user {
		claims[_k] = _v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密算法是HS256时，这里的SignedString必须是[]byte（）类型
	_signedKey := SignedKey
	if len(signedKey) > 0 {
		_signedKey = signedKey[0]
	}
	ss, err := token.SignedString([]byte(_signedKey))
	if err != nil {
		fmt.Println("token生成签名错误,err:", err)
		return ""
	}
	return ss
}

// JwtParseEStoken 解析签名算法为ES256的token
func JwtParseEStoken(tokenES string) string {
	keyX := new(big.Int)
	keyY := new(big.Int)

	keyX.SetString(ECDSAKeyX, 16)
	keyY.SetString(ECDSAKeyY, 16)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     keyX,
		Y:     keyY,
	}

	token, err := jwt.Parse(tokenES, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return &publicKey, nil
	})
	if err != nil {
		fmt.Println("ES256的token解析错误,err:", err)
		return ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["redisValue"].(string)
	}

	fmt.Println("ParseEStoken:Claims类型转换失败")
	return ""
}

// JwtParseHStoken 解析签名算法为HS256的token
func JwtParseHStoken(tokenString string, signedKey ...string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_signedKey := SignedKey
		if len(signedKey) > 0 {
			_signedKey = signedKey[0]
		}
		return []byte(_signedKey), nil
	})
	if err != nil {
		fmt.Println("HS256的token解析错误，err:", err)
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ParseHStoken:claims类型转换失败")
		return nil
	}
	return claims
}

func JwtParseHStoken2(tokenString string, signedKey ...string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_signedKey := SignedKey
		if len(signedKey) > 0 {
			_signedKey = signedKey[0]
		}
		return []byte(_signedKey), nil
	})
	if err != nil {
		fmt.Println("HS256的token解析错误，err:", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ParseHStoken:claims类型转换失败")
		return nil, errors.New("ParseHStoken:claims类型转换失败")
	}
	return claims, nil
}
