// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 21:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : setting.go

package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// GetSettings 获取全部配置
func GetSettings(c *gin.Context) {
	data, resCode := logic.GetSettings()
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, data)
	}
}

// SaveSettings 保存网站信息
func SaveSettings(c *gin.Context) {
	p := new([]model.WebSettings)
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

	// 演示版标签
	if viper.GetString("app.mode") == "demoPro" {
		res.ResErrorWithMsg(c, res.CodeServerBusy, "演示版禁止操作")
		return
	}

	// 处理业务
	resCode := logic.SaveSettings(p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 修改成功
		res.ResSuccess(c, "保存成功")
	}
}
