// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 19:51
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : plugin.go

package controllers

import "C"
import (
	"QLToolsPro/server/cron"
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// PluginData 插件查询
func PluginData(c *gin.Context) {
	// 获取查询类型(ordinary: 普通、cron：定时)
	typeData := c.Query("type")
	if typeData == "ordinary" {
		resCode, data := logic.PluginOrdinaryData()
		switch resCode {
		case res.CodeServerBusy:
			res.ResError(c, res.CodeServerBusy)
		case res.CodeSuccess:
			// 查询成功
			res.ResSuccess(c, data)
		}
	} else {
		resCode, data := logic.PluginCronData()
		switch resCode {
		case res.CodeServerBusy:
			res.ResError(c, res.CodeServerBusy)
		case res.CodeSuccess:
			// 查询成功
			res.ResSuccess(c, data)
		}
	}

}

// PluginCronRefresh 刷新定时插件
func PluginCronRefresh(c *gin.Context) {
	// 暂停已启用任务
	cron.StopTask()

	// 重启定时服务
	if err := cron.Task(); err != nil {
		zap.L().Error("Cron refresh failed, err:" + err.Error())
		res.ResErrorWithMsg(c, res.CodePluginError, "刷新定时插件任务失败")
	} else {
		res.ResSuccess(c, "刷新成功")
	}
}

// PluginUpload 上传插件
func PluginUpload(c *gin.Context) {
	// 获取上传目录(ordinary: 普通、cron：定时)
	typeData := c.Query("type")

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("获取上传文件错误")
		res.ResError(c, res.CodeServerBusy)
		return
	}

	// 获取插件目录绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("获取插件目录绝对路径错误：" + err.Error())
		res.ResError(c, res.CodeServerBusy)
		return
	}

	// 保存文件
	FilePath := ""
	if typeData == "ordinary" {
		FilePath = ExecPath + "/plugin/ordinary/" + file.Filename
	} else {
		FilePath = ExecPath + "/plugin/cron/" + file.Filename
	}

	err = c.SaveUploadedFile(file, FilePath)
	if err != nil {
		zap.L().Error("保存文件错误：" + err.Error())
		res.ResError(c, res.CodeServerBusy)
		return
	}
	rt := "上传成功"
	// 刷新定时插件
	if typeData != "ordinary" {
		go func() {
			// 暂停已启用任务
			cron.StopTask()

			// 重启定时服务
			if err = cron.Task(); err != nil {
				zap.L().Error("Cron refresh failed, err:" + err.Error())
				res.ResErrorWithMsg(c, res.CodePluginError, "刷新定时插件任务失败")
			}
		}()

		rt += rt + ", 定时插件已自动刷新"
	}

	res.ResSuccess(c, rt)
}

// PluginDelete 删除插件
func PluginDelete(c *gin.Context) {
	// 获取参数
	p := new(model.DeletePlugin)
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
	resCode, msg := logic.PluginDelete(p)
	switch resCode {
	case res.CodePluginError:
		// 插件错误
		res.ResErrorWithMsg(c, res.CodePluginError, msg)
	case res.CodeSuccess:
		res.ResSuccess(c, "插件删除成功")
	}
}

/*
普通插件：名字、版本号、作者
定时插件：名字、CRON定时规则、变量名、版本号、作者
前置插件（待定）：名字、版本号、作者
*/

// PluginRemoteDownload 下载远程插件
func PluginRemoteDownload(c *gin.Context) {
	// 获取参数
	p := new(model.DeletePlugin)
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
	resCode, msg := logic.PluginRemoteDownload(p)
	switch resCode {
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodePluginError:
		// 变量错误
		res.ResErrorWithMsg(c, res.CodePluginError, msg)
	case res.CodeSuccess:
		res.ResSuccess(c, "插件下载成功")
	}
}
