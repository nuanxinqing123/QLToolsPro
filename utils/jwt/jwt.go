// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:21
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : jwt.go

package jwt

import (
	"QLToolsPro/server/dao"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

// TokenExpireDuration Token过期时间(7天)
const TokenExpireDuration = time.Hour * 24 * 30

type MyClaims struct {
	UserID     string `json:"user_id"`
	UserSecret string `json:"user_secret"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID, userSecret string) (string, error) {
	// 获取密钥
	jwtKey := dao.GetJWTKey()

	// 加盐
	var mySecret = []byte(jwtKey)
	zap.L().Debug(jwtKey)

	// 创建声明数据
	c := MyClaims{
		UserID:     userID,
		UserSecret: userSecret,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "QLToolsPro",                               // 签发人
		},
	}

	// 使用指定的签名方式创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)

	// 使用指定的签名并获得完整编码后的Token
	return token.SignedString(mySecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 获取密钥
	jwtKey := dao.GetJWTKey()

	// 加盐
	var mySecret = []byte(jwtKey)

	// 解析Token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		zap.L().Warn(err.Error())
		return nil, err
	}

	// 校验Token
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
