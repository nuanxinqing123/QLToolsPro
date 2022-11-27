// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:13
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : code.go

package response

type ResCode int64

const (
	CodeSuccess ResCode = 2000 + iota

	CodeInvalidParam = 5000 + iota
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
	CodeLoginError
	CodeRegisterError
	CodeRePwdError
	CodePanelError
	CodeEnvError
	CodeMessageError
	CodeContainerError
	CodePluginError
	CodeUserEnvError
	CodeCDKError
	CodeNoAdmittance
	CodeAbnormalEnvironment
	CodeAbnormalError
	CodeSystemError
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess: "Success",

	CodeInvalidParam: "请求参数错误",
	CodeServerBusy:   "服务繁忙",

	CodeInvalidToken:        "无效的Token",
	CodeNeedLogin:           "未登录",
	CodeLoginError:          "登录错误",
	CodeRegisterError:       "注册错误",
	CodeRePwdError:          "修改密码错误",
	CodePanelError:          "面板错误",
	CodeEnvError:            "变量错误",
	CodeMessageError:        "推送失败",
	CodeContainerError:      "容器错误",
	CodePluginError:         "插件错误",
	CodeUserEnvError:        "用户变量错误",
	CodeCDKError:            "CD-KEY错误",
	CodeNoAdmittance:        "数据禁止通过",
	CodeAbnormalEnvironment: "环境异常",
	CodeAbnormalError:       "环境错误",
	CodeSystemError:         "系统错误",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
