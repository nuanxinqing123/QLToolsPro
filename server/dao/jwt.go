// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 11:29
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : jwt.go

package dao

import "QLToolsPro/server/model"

// GetJWTKey 获取JWT密钥
func GetJWTKey() string {
	var jwt model.JWTAdmin
	DB.First(&jwt)
	return jwt.SecretKey
}

// CreateJWTKey 创建JWT密钥
func CreateJWTKey(pwd string) {
	var jwt model.JWTAdmin
	jwt.SecretKey = pwd
	DB.Create(&jwt)
}
