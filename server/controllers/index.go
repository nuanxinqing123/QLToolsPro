// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 22:20
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : index.go

package controllers

import (
	"QLToolsPro/server/logic"
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
)

// AdminIndexData 管理员首页数据
func AdminIndexData(c *gin.Context) {
	// 获取查询数据
	resCode, data := logic.AdminIndexData()

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}
