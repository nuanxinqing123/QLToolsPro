package controllers

import (
	"QLToolsPro/server/logic"
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
)

// SystemVersion 系统版本
func SystemVersion(c *gin.Context) {
	data, resCode := logic.CheckVersion()
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, data)
	}
}
