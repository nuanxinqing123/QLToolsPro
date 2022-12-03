// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 10:52
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : container.go

package controllers

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ContainerTransfer 迁移
func ContainerTransfer(c *gin.Context) {
	// 获取参数
	p := new(model.ContainerOperation)
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
	resCode, msg := logic.ContainerTransfer(p)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "操作已进入任务队列, 请稍后前往青龙面板查看结果")
	}
}

// ContainerCopy 容器：复制
func ContainerCopy(c *gin.Context) {
	// 获取参数
	p := new(model.ContainerOperation)
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
	resCode, msg := logic.ContainerCopy(p)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "操作已进入任务队列, 请稍后前往青龙面板查看结果")
	}
}

// ContainerBackup 容器：备份
func ContainerBackup(c *gin.Context) {
	// 获取参数
	p := new(model.ContainerOperationOne)
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
	resCode, msg := logic.ContainerBackup(p)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "操作已进入任务队列, 完成后将自动下载备份文件")
	}
}

// ContainerRestore 容器：恢复
func ContainerRestore(c *gin.Context) {
	// 获取参数
	sID := c.Query("id")
	file, _ := c.FormFile("file")
	zap.L().Debug("【容器：恢复】文件名：" + file.Filename)

	// 保存文件
	err := c.SaveUploadedFile(file, "./"+file.Filename)
	if err != nil {
		// 记录错误
		dao.RecordingError("【容器：恢复】任务", err.Error())
		res.ResError(c, res.CodeServerBusy)
	}

	// 处理业务
	resCode, msg := logic.ContainerRestore(sID)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "操作已进入任务队列, 请稍后前往青龙面板查看结果")
	}
}

// ContainerBackupDownload 容器：备份数据下载
func ContainerBackupDownload(c *gin.Context) {
	Filename := "backup.json"
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", Filename))
	c.File("./" + Filename)
	go logic.DelBackupJSON()
}

// ContainerSynchronization config.sh 同步
func ContainerSynchronization(c *gin.Context) {
	// 获取参数
	p := new(model.ContainerOperationMass)
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
	resCode, msg := logic.ContainerSynchronization(p)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "操作已进入任务队列, 请稍后前往青龙面板查看结果")
	}
}

// ContainerErrorContent 容器：十条日志
func ContainerErrorContent(c *gin.Context) {
	// 处理业务
	resCode, info := logic.ContainerErrorContent()
	switch resCode {
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, info)
	}
}

// ContainerCronBackup 定时备份面板变量
func ContainerCronBackup(c *gin.Context) {
	// 获取参数
	p := new(model.CronBackUpEnv)
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
	resCode, msg := logic.ContainerCronBackup(p)
	switch resCode {
	case res.CodeContainerError:
		res.ResErrorWithMsg(c, res.CodeContainerError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		res.ResSuccess(c, "定时任务已设置成功")
	}
}
