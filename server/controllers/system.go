package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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

// SystemSoftwareUpdate 更新系统
func SystemSoftwareUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.SoftWareGOOS)
	if err := c.ShouldBindJSON(&p); err != nil {
		// 参数校验
		zap.L().Error("SignInHandle with invalid param", zap.Error(err))

		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			res.ResError(c, res.CodeInvalidParam)
			return
		}

		// 翻译错误
		res.ResErrorWithMsg(c, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
		return
	}

	resCode, msg := logic.SystemSoftwareUpdate(p)
	switch resCode {
	case res.CodeSystemError:
		res.ResErrorWithMsg(c, res.CodeSystemError, msg)
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, msg)
	}
}
