// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 10:51
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : user.go

package controllers

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/logic"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	val "QLToolsPro/utils/validator"
	"QLToolsPro/utils/wxpusher"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const CtxUserIDKey = "userID"

// CreateVerificationCode 创建验证码
func CreateVerificationCode(c *gin.Context) {
	id, bs64, err := logic.CaptMake()
	if err != nil {
		// 服务繁忙
		res.ResError(c, res.CodeServerBusy)
	} else {
		res.ResSuccess(c, gin.H{
			"id":   id,
			"bs64": bs64,
		})
	}
}

// SignUpHandle 注册请求
func SignUpHandle(c *gin.Context) {
	// 获取参数
	p := new(model.UserSignUp)
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
	resCode, msg := logic.SignUp(p)
	switch resCode {
	case res.CodeRegisterError:
		// 注册错误
		res.ResErrorWithMsg(c, res.CodeRegisterError, msg)
	case res.CodeServerBusy:
		// 内部服务错误
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 注册成功
		res.ResSuccess(c, "注册完成")
	}
}

// SignInHandle 登录请求
func SignInHandle(c *gin.Context) {
	// 获取参数
	p := new(model.UserSignIn)
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
	var RemoteIP string
	// 获取IP地址
	if "127.0.0.1" == c.RemoteIP() {
		RemoteIP = c.GetHeader("X-Real-IP")
	} else {
		RemoteIP = c.RemoteIP()
	}

	resCode, msg := logic.SignIn(p, RemoteIP)
	switch resCode {
	case res.CodeAbnormalEnvironment:
		// 登录环境异常
		res.ResErrorWithMsg(c, res.CodeAbnormalEnvironment, msg)
	case res.CodeLoginError:
		// 邮箱或者密码错误
		res.ResErrorWithMsg(c, res.CodeLoginError, msg)
	case res.CodeServerBusy:
		// 生成Token出错
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 登录成功,返回Token
		res.ResSuccess(c, msg)
	}
}

// AppletLoginHandle 小程序登录
func AppletLoginHandle(c *gin.Context) {
	// 获取参数
	p := new(model.AppletLogin)
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
	resCode, msg := logic.AppletLogin(p)
	switch resCode {
	case res.CodeLoginError:
		// 注册错误
		res.ResErrorWithMsg(c, res.CodeLoginError, msg)
	case res.CodeServerBusy:
		// 内部服务错误
		res.ResError(c, res.CodeServerBusy)
	case res.CodeAbnormalEnvironment:
		// 登录环境异常
		res.ResErrorWithMsg(c, res.CodeAbnormalEnvironment, msg)
	case res.CodeSuccess:
		// 登录完成
		var RemoteIP string
		UID, _ := c.Get(CtxUserIDKey)
		// 获取IP地址
		if "127.0.0.1" == c.RemoteIP() {
			RemoteIP = c.GetHeader("X-Real-IP")
		} else {
			RemoteIP = c.RemoteIP()
		}
		go dao.UpdateUserLoginIP(RemoteIP, UID)
		res.ResSuccess(c, msg)
	}
}

// GetUserOneData 用户信息：获取
func GetUserOneData(c *gin.Context) {
	// 处理业务
	UID, _ := c.Get(CtxUserIDKey)
	resCode, data := logic.GetUserOneData(UID)
	switch resCode {
	case res.CodeSuccess:
		// 更新成功
		res.ResSuccess(c, data)
	}
}

// UserSettings 获取单个配置
func UserSettings(c *gin.Context) {
	key := c.Query("key")
	data, resCode := logic.UserSettings(key)
	switch resCode {
	case res.CodeServerBusy:
		// 越权
		res.ResErrorWithMsg(c, res.CodeServerBusy, "获取内容为空")
	case res.CodeSuccess:
		// 获取成功
		res.ResSuccess(c, data)
	}
}

// UserWxpusherState 用户：检查用户WxPusher订阅状态
func UserWxpusherState(c *gin.Context) {
	UID, _ := c.Get(CtxUserIDKey)
	wxpusherID := dao.GetUserWxpusherID(UID)
	if wxpusherID == "" {
		// 未订阅
		res.ResSuccess(c, false)
	} else {
		// 已订阅
		res.ResSuccess(c, true)
	}
}

// UserWxpusherQrcode 用户：获取WxPusher订阅二维码
func UserWxpusherQrcode(c *gin.Context) {
	data, err := wxpusher.GetWxPusherQRCode()
	if err != nil {
		zap.L().Error("[获取WxPusher订阅二维码]错误，原因：" + err.Error())
		res.ResError(c, res.CodeServerBusy)
		return
	}

	// 获取成功
	res.ResSuccess(c, data)
}

// UserWxpusherUpdate 用户：WxPusher更新
func UserWxpusherUpdate(c *gin.Context) {
	wxUid := c.Query("wx_uid")
	UID, _ := c.Get(CtxUserIDKey)
	err := dao.UpdateUserWxpusherID(UID, wxUid)
	if err != nil {
		// 未订阅
		res.ResError(c, res.CodeServerBusy)
	} else {
		// 已订阅
		res.ResSuccess(c, "更新成功")
	}
}

// UserDivisionData 用户分页查询
func UserDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")

	resCode, data := logic.UserDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// UserSearch 用户筛选搜索
func UserSearch(c *gin.Context) {
	// 用户名模糊查询
	s := c.Query("s")
	resCode, data := logic.UserSearch(s)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// UserInformationUpdate 用户数据更新
func UserInformationUpdate(c *gin.Context) {
	// 获取参数
	p := new(model.UpdateUserData)
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
	resCode := logic.UserInformationUpdate(p)
	switch resCode {
	case res.CodeSuccess:
		// 更新成功
		res.ResSuccess(c, "更新成功")
	}
}

// UserInformationDelete 用户数据删除
func UserInformationDelete(c *gin.Context) {
	// 获取参数
	p := new(model.DeleteUserData)
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
	resCode := logic.UserInformationDelete(p)
	switch resCode {
	case res.CodeSuccess:
		// 删除成功
		res.ResSuccess(c, "删除成功")
	}
}

// UserFindPwd 找回密码 - 发送Token
func UserFindPwd(c *gin.Context) {
	// 获取参数
	p := new(model.UserFindPwd)
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
	resCode, msg := logic.UserFindPwd(p)
	switch resCode {
	case res.CodeRePwdError:
		// 修改密码错误
		res.ResErrorWithMsg(c, res.CodeRePwdError, msg)
	case res.CodeSuccess:
		// 修改成功
		res.ResSuccess(c, "发送成功")
	}
}

// UserRePwd 找回密码 - 修改密码
func UserRePwd(c *gin.Context) {
	// 获取参数
	p := new(model.UserRePwd)
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
	resCode, msg := logic.UserRePwd(p)
	switch resCode {
	case res.CodeRePwdError:
		// 修改密码错误
		res.ResErrorWithMsg(c, res.CodeRePwdError, msg)
	case res.CodeSuccess:
		// 修改成功
		var RemoteIP string
		UID, _ := c.Get(CtxUserIDKey)
		// 获取IP地址
		if "127.0.0.1" == c.RemoteIP() {
			RemoteIP = c.GetHeader("X-Real-IP")
		} else {
			RemoteIP = c.RemoteIP()
		}
		go dao.UpdateUserLoginIP(RemoteIP, UID)
		res.ResSuccess(c, "修改成功")
	}
}

// UserAbnormalCode 登录异常-发送验证码
func UserAbnormalCode(c *gin.Context) {
	// 获取参数
	p := new(model.UserAbnormalEmail)
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
	resCode, msg := logic.AbnormalEmail(p)
	switch resCode {
	case res.CodeAbnormalError:
		res.ResErrorWithMsg(c, res.CodeAbnormalError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 发送成功
		res.ResSuccess(c, "发送成功，请前往微信公众号【WxPusher】查看验证码")
	}
}

// UserAbnormalSignin 登录异常-登录
func UserAbnormalSignin(c *gin.Context) {
	// 获取参数
	p := new(model.UserAbnormalSignin)
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
	var RemoteIP string
	// 获取IP地址
	if "127.0.0.1" == c.RemoteIP() {
		RemoteIP = c.GetHeader("X-Real-IP")
	} else {
		RemoteIP = c.RemoteIP()
	}

	resCode, msg := logic.AbnormalSignin(p, RemoteIP)
	switch resCode {
	case res.CodeAbnormalError:
		// 修改密码错误
		res.ResErrorWithMsg(c, res.CodeAbnormalError, msg)
	case res.CodeServerBusy:
		res.ResError(c, res.CodeServerBusy)
	case res.CodeSuccess:
		// 登录成功
		res.ResSuccess(c, msg)
	}
}
