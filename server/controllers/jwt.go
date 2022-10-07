// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 11:40
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : jwt.go

package controllers

import (
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
)

// CheckToken 检查Token是否有效
func CheckToken(c *gin.Context) {
	res.ResSuccess(c, "有效Token")
}
