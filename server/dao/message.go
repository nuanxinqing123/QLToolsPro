// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 18:16
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : message.go

package dao

import "QLToolsPro/server/model"

// GetDivisionMessageData 条件查询面板数据
func GetDivisionMessageData(page, q int) []model.Message {
	var message []model.Message
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Model(&model.User{}).Select("user_id", "username", "user_wxpusher").Order("id desc").Limit(25).Offset(0).Find(&message)
	} else {
		// 获取第N页数据
		DB.Model(&model.User{}).Select("user_id", "username", "user_wxpusher").Order("id desc").Limit(25).Offset((page - 1) * q).Find(&message)
	}
	return message
}

// GetMessageDataPage 获取Message表总数据
func GetMessageDataPage() int64 {
	var c []model.User
	result := DB.Where("user_wxpusher <> ''").Find(&c)
	return result.RowsAffected
}

// GetBindWxPusherUserData 获取绑定WxPusherID的用户信息
func GetBindWxPusherUserData() []model.User {
	var c []model.User
	DB.Where("user_wxpusher <> ''").Find(&c)
	return c
}
