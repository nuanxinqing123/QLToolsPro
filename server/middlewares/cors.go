// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:09
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cors.go

package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 配置跨域
func Cors() gin.HandlerFunc {
	//return func(c *gin.Context) {
	//	method := c.Request.Method
	//
	//	c.Header("Access-Control-Allow-Origin", "*")
	//	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	//	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	//	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	//	c.Header("Access-Control-Allow-Credentials", "true")
	//
	//	// 放行所有OPTIONS方法
	//	if method == "OPTIONS" {
	//		c.AbortWithStatus(http.StatusNoContent)
	//	}
	//	// 处理请求
	//	c.Next()
	//}
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,DELETE,OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
