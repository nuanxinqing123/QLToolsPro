// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 18:16
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : message.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"QLToolsPro/utils/wxpusher"
	"go.uber.org/zap"
	"strconv"
)

// MessageDivisionData 面板分页查询
func MessageDivisionData(page, quantity string) (res.ResCode, model.MessagePageData) {
	var data []model.Message
	var pageData model.MessagePageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.GetDivisionMessageData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.GetDivisionMessageData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.GetDivisionMessageData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetMessageDataPage()
	// 计算页数
	z := count / int64(q)
	var y int64
	y = count % int64(q)

	if y != 0 {
		pageData.Page = z + 1
	} else {
		pageData.Page = z
	}
	pageData.PageData = data

	return res.CodeSuccess, pageData
}

// MessageSend 管理员消息群发
func MessageSend(p *model.AdminMessage) (res.ResCode, string) {
	b, msg := wxpusher.AdminSendMessage(p.UserWxpusher, p.Msg)
	if b {
		return res.CodeSuccess, msg
	}

	return res.CodeMessageError, msg
}

// MessageSendAll 管理员全体消息发送
func MessageSendAll(p *model.AdminMessageAll) (res.ResCode, string) {
	// 查询所有以绑定WxPusherID的用户
	var UserList []string

	user := dao.GetBindWxPusherUserData()
	for _, u := range user {
		UserList = append(UserList, u.UserWxpusher)
	}

	b, msg := wxpusher.AdminSendMessage(UserList, p.Msg)
	if b {
		return res.CodeSuccess, msg
	}

	return res.CodeMessageError, msg
}
