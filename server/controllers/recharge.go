// -*- coding: utf-8 -*-
// @Time    : 2022/8/31 10:25
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : recharge.go

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

// UserRechargeIntegral 用户充值：用户积分充值
func UserRechargeIntegral(c *gin.Context) {
	// 获取参数
	p := new(model.UserRecharge)
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

	// 处理业务
	UID, _ := c.Get(CtxUserIDKey)
	resCode, msg := logic.UserRechargeIntegral(UID, p)
	switch resCode {
	case res.CodeCDKError:
		// CDK错误
		res.ResErrorWithMsg(c, res.CodeCDKError, msg)
	case res.CodeSuccess:
		// 充值成功
		res.ResSuccess(c, "充值成功")
	}
}

// UserRechargeVIP 用户充值：用户会员充值
func UserRechargeVIP(c *gin.Context) {
	// 获取参数
	p := new(model.UserRecharge)
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

	// 处理业务
	UID, _ := c.Get(CtxUserIDKey)
	resCode, msg := logic.UserRechargeVIP(UID, p)
	switch resCode {
	case res.CodeCDKError:
		// CDK错误
		res.ResErrorWithMsg(c, res.CodeCDKError, msg)
	case res.CodeSuccess:
		// 充值成功
		res.ResSuccess(c, "充值成功")
	}
}

// RechargeUserDivisionData 用户充值：记录查询
func RechargeUserDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	// 处理业务
	UID, _ := c.Get(CtxUserIDKey)
	resCode, msg := logic.RechargeUserDivisionData(UID, page, quantity)
	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, msg)
	}
}

// RechargeDivisionData 充值数据：以20条数据分割
func RechargeDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")
	resCode, data := logic.RechargeDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// RechargeSearch 充值数据：CDK搜索/UserID搜索
func RechargeSearch(c *gin.Context) {
	var data []model.Recharge
	var resCode res.ResCode

	s := c.Query("s")
	if s == "" {
		res.ResErrorWithMsg(c, res.CodeInvalidParam, "请求数据不完整")
		return
	}
	resCode, data = logic.RechargeSearch(s)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// CDKEYUserRechargeIntegral 指定用户充值
func CDKEYUserRechargeIntegral(c *gin.Context) {
	// 获取参数
	p := new(model.AdminRecharge)
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

	// 处理业务
	resCode, msg := logic.CDKEYUserRechargeIntegral(p)
	switch resCode {
	case res.CodeCDKError:
		// CDK错误
		res.ResErrorWithMsg(c, res.CodeCDKError, msg)
	case res.CodeSuccess:
		// 充值成功
		res.ResSuccess(c, "充值成功")
	}
}
