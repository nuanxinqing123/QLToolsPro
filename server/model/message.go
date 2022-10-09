// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 18:16
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : message.go

package model

// Message 已绑定WxPusher的用户信息
type Message struct {
	UserID       string // 用户ID
	Username     string // 用户名
	UserWxpusher string // 用户Wxpusher ID
}

// MessagePageData 用户WxPusher分页查询
type MessagePageData struct {
	Page     int64     `json:"page"`
	PageData []Message `json:"pageData"`
}

// AdminMessage 管理员消息群发
type AdminMessage struct {
	UserWxpusher []string `json:"user_wxpusher" binding:"required"`
	Msg          string   `json:"message" binding:"required"`
}

// AdminMessageAll 管理员全体消息发送
type AdminMessageAll struct {
	Msg string `json:"message" binding:"required"`
}
