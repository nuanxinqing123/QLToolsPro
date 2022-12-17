// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:15
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : env.go

package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// EnvDivisionData 变量分页查询
func EnvDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	resCode, data := logic.EnvDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// EnvDivisionDataSimple 获取变量简易数据
func EnvDivisionDataSimple(c *gin.Context) {
	resCode, data := logic.EnvDivisionDataSimple()

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// EnvAdd 变量新增
func EnvAdd(c *gin.Context) {
	// 获取参数
	p := new(model.EnvAdd)
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
	resCode, msg := logic.EnvAdd(p)
	switch resCode {
	case res.CodeEnvError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodeEnvError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "变量信息创建成功")
	}
}

// EnvUpdate 变量修改
func EnvUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.EnvUpdate)
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
	zap.L().Debug(fmt.Sprintln(p))
	resCode, msg := logic.EnvUpdate(p)
	switch resCode {
	case res.CodeEnvError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodeEnvError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "变量信息更新成功")
	}
}

// EnvDelete 变量删除
func EnvDelete(c *gin.Context) {
	// 获取参数
	p := new(model.EnvDelete)
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
	resCode, msg := logic.EnvDelete(p)
	switch resCode {
	case res.CodeEnvError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodeEnvError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "变量信息删除成功")
	}
}

// UserEnvDataSearch 用户变量管理:筛选查询
func UserEnvDataSearch(c *gin.Context) {
	// 查询类型（面板查询：panel、用户查询：user）
	typeData := c.Query("type")
	search := c.Query("s")
	zap.L().Debug("[用户变量管理:筛选查询]s：" + search)
	resCode, data := logic.UserEnvDataSearch(typeData, search)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// UserEnvDataUpdate 用户变量管理:修改用户变量
func UserEnvDataUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.UserEnvUpdate)
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
	resCode, msg := logic.UserEnvDataUpdate(p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeUserEnvError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodeUserEnvError, msg)
	case res.CodeSuccess:
		res.ResSuccess(c, "变量信息更新成功")
	}
}

// UserEnvDataDelete 用户变量管理:删除用户变量
func UserEnvDataDelete(c *gin.Context) {
	// 获取参数
	p := new(model.PanelEnvDelete)
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
	resCode, msg := logic.UserEnvDataDelete(p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeUserEnvError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodeUserEnvError, msg)
	case res.CodeSuccess:
		go panel.UpdateOnlineServerCache()
		res.ResSuccess(c, "变量信息删除成功")
	}
}

// UserRefresh 手动刷新首页缓存数据
func UserRefresh(c *gin.Context) {
	go panel.UpdateOnlineServerCache()
	res.ResSuccess(c, "缓存刷新成功")
}
