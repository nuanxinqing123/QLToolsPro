// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 18:16
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : message.go

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

// MessageDivisionData 用户WxPusher分页查询
func MessageDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	resCode, data := logic.MessageDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// MessageSend 管理员消息群发
func MessageSend(c *gin.Context) {
	// 获取参数
	p := new(model.AdminMessage)
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

	// 处理业务
	resCode, msg := logic.MessageSend(p)
	switch resCode {
	case res.CodeMessageError:
		// 推送失败
		res.ResErrorWithMsg(c, res.CodeMessageError, msg)
	case res.CodeSuccess:
		res.ResSuccess(c, "消息群发成功")
	}
}

// MessageSendAll 管理员全体消息发送
func MessageSendAll(c *gin.Context) {
	// 获取参数
	p := new(model.AdminMessageAll)
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

	// 处理业务
	resCode, msg := logic.MessageSendAll(p)
	switch resCode {
	case res.CodeMessageError:
		// 推送失败
		res.ResErrorWithMsg(c, res.CodeMessageError, msg)
	case res.CodeSuccess:
		res.ResSuccess(c, "消息群发成功")
	}
}
