// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:09
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : auth.go

package middlewares

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/utils/jwt"
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
	"strings"
)

const CtxUserIDKey = "userID"

// UserAuth 用户基于JWT的认证中间件
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			res.ResError(c, res.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.ResError(c, res.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			res.ResError(c, res.CodeInvalidToken)
			c.Abort()
			return
		}

		// 校验是否属于注册用户
		cUser := dao.CheckIsUser(mc.UserID)

		if cUser.Username == "" && cUser.Email == "" && cUser.IsState == false {
			c.Abort()
			res.ResErrorWithMsg(c, res.CodeInvalidToken, "无访问权限或认证已过期")
			return
		} else {
			if cUser.IsState == false {
				// 账号已被封禁
				res.ResError(c, res.CodeInvalidToken)
				c.Abort()
				return
			}

			if mc.UserSecret != cUser.Password[:6] {
				// 已修改密码，需要强制下线
				res.ResError(c, res.CodeInvalidToken)
				c.Abort()
				return
			} else {
				//将当前请求的userID信息保存到请求的上下文c上
				c.Set(CtxUserIDKey, mc.UserID)
				c.Next()
			}
		}
	}
}

// AdminAuth 管理员JWT的认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			res.ResError(c, res.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			res.ResError(c, res.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			res.ResError(c, res.CodeInvalidToken)
			c.Abort()
			return
		}

		// 检查是否属于管理员
		cAdmin := dao.CheckIsUser(mc.UserID)

		if cAdmin.IsAdmin == false {
			c.Abort()
			res.ResErrorWithMsg(c, res.CodeInvalidToken, "无访问权限或认证已过期")
			return
		} else {
			if cAdmin.IsState == false {
				// 账号已被封禁
				res.ResError(c, res.CodeInvalidToken)
				c.Abort()
				return
			}

			if mc.UserSecret != cAdmin.Password[:6] == true {
				// 已修改密码，需要强制下线
				res.ResError(c, res.CodeInvalidToken)
				c.Abort()
				return
			} else {
				//将当前请求的userID信息保存到请求的上下文c上
				c.Set(CtxUserIDKey, mc.UserID)
				c.Next()
			}
		}
	}
}
