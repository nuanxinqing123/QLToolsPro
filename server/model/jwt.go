// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 11:29
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : jwt.go

package model

import "gorm.io/gorm"

// JWTAdmin JWT表
type JWTAdmin struct {
	gorm.Model
	SecretKey string // JWT签名密钥
}

// CheckToken 检查Token是否有效
type CheckToken struct {
	JWToken string `json:"token" binding:"required"`
}
