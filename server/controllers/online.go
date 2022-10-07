// -*- coding: utf-8 -*-
// @Time    : 2022/9/3 11:58
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : online.go

package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// OnlineServer 在线服务
func OnlineServer(c *gin.Context) {
	// 处理业务
	resCode, data := logic.OnlineServer()
	switch resCode {
	case res.CodeServerBusy:
		// 内部服务错误
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 登录成功,返回Token
		res.ResSuccess(c, data)
	}
}

// OnlineUploadData 上传变量
func OnlineUploadData(c *gin.Context) {
	// 获取参数
	p := new(model.OnlineEnvUpload)
	if err := c.ShouldBindJSON(&p); err != nil {
		// 参数校验
		zap.L().Error("SignUpHandle with invalid param", zap.Error(err))

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

	// 业务处理
	UID, _ := c.Get(CtxUserIDKey)
	resCode, msg := logic.OnlineUploadData(UID, p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResErrorWithMsg(c, res.CodeServerBusy, "服务繁忙,请稍后重试")
	case res.CodeEnvError:
		res.ResErrorWithMsg(c, res.CodeEnvError, msg)
	case res.CodeNoAdmittance:
		res.ResErrorWithMsg(c, res.CodeNoAdmittance, msg)
	case res.CodeSuccess:
		// 上传成功
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "上传成功")
	}
}
