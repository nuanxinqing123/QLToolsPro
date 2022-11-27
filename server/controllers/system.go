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

// SystemSoftwareUpdate 更新系统
//func SystemSoftwareUpdate(c *gin.Context) {
//	// 获取参数
//	p := new(model.SoftWareGOOS)
//	if err := c.ShouldBindJSON(&p); err != nil {
//		// 参数校验
//		zap.L().Error("SignInHandle with invalid param", zap.Error(err))
//
//		// 判断err是不是validator.ValidationErrors类型
//		errs, ok := err.(validator.ValidationErrors)
//		if !ok {
//			res.ResError(c, res.CodeInvalidParam)
//			return
//		}
//
//		// 翻译错误
//		res.ResErrorWithMsg(c, res.CodeInvalidParam, val.RemoveTopStruct(errs.Translate(val.Trans)))
//		return
//	}
//
//	resCode, txt := logic.SystemSoftwareUpdate(p)
//	switch resCode {
//	case res.CodeUpdateServerBusy:
//		res.ResErrorWithMsg(c, res.CodeUpdateServerBusy, txt)
//	case res.CodeSuccess:
//		// 获取成功
//		res.ResSuccess(c, txt)
//	}
//}
