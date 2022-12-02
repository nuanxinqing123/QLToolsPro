// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 10:54
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : user.go

package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	/*
		gorm.Model：基础结构（ID、CreatedAt、UpdatedAt、DeletedAt）
	*/
	gorm.Model
	UserID         string    // 用户ID
	Email          string    // 用户邮箱
	AppOpenID      string    // 小程序用户唯一标识
	Username       string    // 用户名
	Password       string    // 用户密码
	Integral       int64     // 用户积分
	ActivationTime time.Time // 会员到期时间
	IsAdmin        bool      // 是否属于管理员（true：是、false：否）
	UserWxpusher   string    // 用户Wxpusher ID
	IsState        bool      // 用户状态（true：启用、false：封禁）
	LoginIP        string    // 用户上次登录IP
}

// UserData 用户信息
type UserData struct {
	UserID           string // 用户ID
	Email            string // 用户邮箱
	Username         string // 用户名
	Integral         int64  // 用户积分
	UserWxpusher     string // 用户Wxpusher ID
	IsVIP            bool   // 是否属于VIP用户（true：是、false：否）
	IsActivationTime bool   // 会员是否到期
	ActivationTime   string // 会员到期时间
	IsAdmin          bool   // 是否属于管理员（true：是、false：否）
}

// UserSignUp 用户注册
type UserSignUp struct {
	Email      string `json:"email" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Capt       string `json:"capt" binding:"required"`
	Id         string `json:"id" binding:"required"`
}

// UserSignIn 用户登录
type UserSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Capt     string `json:"capt" binding:"required"`
	Id       string `json:"id" binding:"required"`
}

// AppletLogin 小程序登录
type AppletLogin struct {
	// 用户标识Code
	Code string `json:"code" binding:"required"`
}

// AppletRes 小程序获取返回
type AppletRes struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

// UserPageData 用户分页数据
type UserPageData struct {
	Page     int64  `json:"page"`
	PageData []User `json:"pageData"`
}

// UpdateUserData 用户数据修改
type UpdateUserData struct {
	ID           int64  `json:"id"  binding:"required"`
	UserWxpusher string `json:"user_wxpusher"` // 用户WxpusherID
	IsAdmin      bool   `json:"is_admin"`      // 是否属于管理员（true：是、false：否）
	IsState      bool   `json:"is_state"`      // 用户状态（true：启用、false：封禁）
}

// DeleteUserData 用户数据删除
type DeleteUserData struct {
	ID int `json:"id"  binding:"required"`
}

// UserFindPwd 用户找回密码 - 发送Token
type UserFindPwd struct {
	FindType string `json:"find_type" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Capt     string `json:"capt" binding:"required"`
	Id       string `json:"id" binding:"required"`
}

// CacheRecordPwd 缓存找回密码数据
type CacheRecordPwd struct {
	Code   string
	UserId string
}

// UserRePwd 用户找回密码 - 修改密码
type UserRePwd struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserAbnormalEmail 登录异常 - 发送验证码
type UserAbnormalEmail struct {
	UserName string `json:"username" binding:"required"`
}

// UserAbnormalSignin 登录异常 - 登录
type UserAbnormalSignin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	VfCode   string `json:"id" binding:"required"`
}
