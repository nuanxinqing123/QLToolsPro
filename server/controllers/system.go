package controllers

import (
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os/exec"
	"strconv"
	"syscall"
	"time"
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
func SystemSoftwareUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.SoftWareGOOS)
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

	resCode, msg := logic.SystemSoftwareUpdate(p)
	switch resCode {
	case res.CodeSystemError:
		res.ResErrorWithMsg(c, res.CodeSystemError, msg)
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, msg)
	}
}

// SystemState 关闭/重启系统
func SystemState(c *gin.Context) {
	// shutdown：关闭、restart：重启
	t := c.Query("type")
	if t == "shutdown" {
		go func() {
			time.Sleep(time.Second * 3)
			zap.L().Debug("进程PID：" + strconv.Itoa(syscall.Getpid()))
			cmd := exec.Command("/bin/bash", "-c", "kill -SIGTERM "+strconv.Itoa(syscall.Getpid()))
			err := cmd.Start()
			if err != nil {
				zap.L().Error("[关闭]：" + err.Error())
			}
			err = cmd.Start()
		}()

	} else {
		//go func() {
		//	time.Sleep(time.Second * 3)
		//	zap.L().Debug("进程PID：" + strconv.Itoa(syscall.Getpid()))
		//	cmd := exec.Command("/bin/bash", "-c", "kill -SIGHUP "+strconv.Itoa(syscall.Getpid()))
		//
		//	// 调用命令，如果发生错误，记录错误信息
		//	err := cmd.Start()
		//	if err != nil {
		//		zap.L().Error("[重启]：" + err.Error())
		//	}
		//}()
		res.ResSuccess(c, "系统重启功能暂时无法操作")
		return
	}
	res.ResSuccess(c, "系统将在三秒后执行操作")
}

// TaskDataQuery 查询任务
func TaskDataQuery(c *gin.Context) {
	tp := c.Query("type")

	data := logic.TaskDataQuery(tp)

	// 查询成功
	res.ResSuccess(c, data)
}
