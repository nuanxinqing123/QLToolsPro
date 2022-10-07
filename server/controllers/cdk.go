// -*- coding: utf-8 -*-
// @Time    : 2022/8/26 10:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cdk.go

package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CDKEYDivisionCDKData CD-KEY分页查询
func CDKEYDivisionCDKData(c *gin.Context) {
	// 获取查询页码
	ctype := c.Query("ctype")
	page := c.Query("page")
	quantity := c.Query("quantity")
	resCode, data := logic.CDKEYDivisionCDKData(ctype, page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// CDKEYSearch CD-KEY数据查询
func CDKEYSearch(c *gin.Context) {
	ctype := c.Query("ctype")
	s := c.Query("s")
	zap.L().Debug("【CD-KEY搜索】值：" + s)
	if s == "" {
		res.ResErrorWithMsg(c, res.CodeInvalidParam, "请求数据不完整")
		return
	}
	resCode, data := logic.CDKEYSearch(ctype, s)

	switch resCode {
	case res.CodeCDKError:
		res.ResErrorWithMsg(c, res.CodeCDKError, "查询数据不存在")
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// CDKEYRemarksSearch CD-KEY标识查询
func CDKEYRemarksSearch(c *gin.Context) {
	// 查询内容
	s := c.Query("s")

	resCode, data := logic.CDKEYRemarksSearch(s)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// CDKEYAdd 批量生成CD-KEY
func CDKEYAdd(c *gin.Context) {
	// 获取参数
	p := new(model.CreateCDK)
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
	resCode := logic.CDKEYAdd(p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResErrorWithMsg(c, res.CodeServerBusy, "生成CDK失败，请检查日志获取报错信息")
	case res.CodeCDKError:
		res.ResErrorWithMsg(c, res.CodeCDKError, "创建CDK已写入数据库，但生成下载文件失败")
	case res.CodeSuccess:
		// 生成成功
		res.ResSuccess(c, "生成成功")
	}
}

// CDKEYDataDownload 下载CD-KEY文件
func CDKEYDataDownload(c *gin.Context) {
	Filename := "CDK.txt"
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", Filename))
	c.File("./" + Filename)
	go logic.CDKEYDataDelete()
}

// CDKEYBatchUpdate 修改CD-KEY
func CDKEYBatchUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.UpdateCDK)
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
	resCode := logic.CDKEYBatchUpdate(p)
	switch resCode {
	case res.CodeEnvError:
		res.ResErrorWithMsg(c, res.CodeEnvError, "更新CD-KEY失败")
	case res.CodeSuccess:
		// 更新成功
		res.ResSuccess(c, "更新成功")
	}
}

// CDKEYDelete 删除CD-KEY数据
func CDKEYDelete(c *gin.Context) {
	// 获取参数
	p := new(model.DelCDK)
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
	resCode := logic.CDKEYDelete(p)
	switch resCode {
	case res.CodeEnvError:
		res.ResErrorWithMsg(c, res.CodeEnvError, "删除CD-KEY失败")
	case res.CodeSuccess:
		// 删除成功
		res.ResSuccess(c, "删除成功")
	}
}
