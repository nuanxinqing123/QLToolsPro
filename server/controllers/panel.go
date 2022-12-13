// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:33
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : panel.go

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

// PanelDivisionData 面板分页查询
func PanelDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	resCode, data := logic.PanelDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// PanelDivisionDataSimple 获取面板简易数据
func PanelDivisionDataSimple(c *gin.Context) {
	resCode, data := logic.PanelDivisionDataSimple()

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// PanelAdd 面板新增
func PanelAdd(c *gin.Context) {
	// 获取参数
	p := new(model.PanelAdd)
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
	resCode, msg := logic.PanelAdd(p)
	switch resCode {
	case res.CodePanelError:
		// 面板错误
		res.ResErrorWithMsg(c, res.CodePanelError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "面板信息创建成功")
	}
}

// PanelUpdate 面板修改
func PanelUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.PanelUpdate)
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
	resCode, msg := logic.PanelUpdate(p)
	switch resCode {
	case res.CodePanelError:
		// 面板错误
		res.ResErrorWithMsg(c, res.CodePanelError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "面板信息更新成功")
	}
}

// PanelDelete 面板删除
func PanelDelete(c *gin.Context) {
	// 获取参数
	p := new(model.PanelDelete)
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
	resCode, msg := logic.PanelDelete(p)
	switch resCode {
	case res.CodePanelError:
		// 面板错误
		res.ResErrorWithMsg(c, res.CodePanelError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "面板信息删除成功")
	}
}

// PanelTestConnect 面板测试连接
func PanelTestConnect(c *gin.Context) {
	// 获取参数
	p := new(model.PanelAdd)
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
	resCode, token := panel.TestGetPanelToken(p.URL, p.ID, p.Secret)
	if resCode == res.CodeServerBusy {
		// 内部服务错误
		res.ResError(c, res.CodeServerBusy)
		return
	}

	if token.Code != 200 {
		// 授权错误
		res.ResErrorWithMsg(c, res.CodePanelError, "client_id或client_secret有误")
		return
	}

	res.ResSuccess(c, "面板连接测试成功")
	return
}

// PanelBindingUpdate 面板修改绑定变量
func PanelBindingUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.PanelBindingUpdate)
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
	resCode, msg := logic.PanelBindingUpdate(p)
	switch resCode {
	case res.CodePanelError:
		// 面板错误
		res.ResErrorWithMsg(c, res.CodePanelError, msg)
	case res.CodeSuccess:
		// 修改成功
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "修改成功")
	}
}

// PanelUnbindUpdate 面板解除变量绑定
func PanelUnbindUpdate(c *gin.Context) {
	// 处理业务
	resCode := logic.PanelUnbindUpdate()
	switch resCode {
	case res.CodeSuccess:
		// 解绑成功
		res.ResSuccess(c, "解绑成功")
	}
}

// PanelTokenUpdate 面板批量更新Token
func PanelTokenUpdate(c *gin.Context) {
	// 处理业务
	resCode := logic.PanelTokenUpdate()
	switch resCode {
	case res.CodeSuccess:
		// 更新成功
		res.ResSuccess(c, "面板批量更新Token成功")
	}
}
