// -*- coding: utf-8 -*-
// @Time    : 2022/8/27 11:22
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : record.go

package controllers

import (
	"QLToolsPro/server/logic"
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
)

// RecordDivisionData 上传记录分页查询
func RecordDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")
	resCode, data := logic.RecordDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// RecordSearch 上传记录数据查询
func RecordSearch(c *gin.Context) {
	s := c.Query("s")
	if s == "" {
		res.ResErrorWithMsg(c, res.CodeInvalidParam, "请求数据不完整")
		return
	}
	resCode, data := logic.RecordSearch(s)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// RecordUserDivisionData 用户上传：记录查询
func RecordUserDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	// 处理业务
	UID, _ := c.Get(CtxUserIDKey)
	resCode, msg := logic.RecordUserDivisionData(UID, page, quantity)
	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, msg)
	}
}
