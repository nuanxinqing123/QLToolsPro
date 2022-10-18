// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 23:17
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : index.go

package dao

import "QLToolsPro/server/model"

// GetPanelCount 获取面板存在数量
func GetPanelCount() int64 {
	var p model.Panel
	return DB.Find(&p).RowsAffected
}

// GetEnvCount 获取变量存在数量
func GetEnvCount() int64 {
	return DB.Find(&model.Env{}).RowsAffected
}

// GetRegisterUserCount 获取网站注册用户数量
func GetRegisterUserCount() int64 {
	return DB.Find(&model.User{}).RowsAffected
}

// GetVIPUserCount 获取网站会员用户数量
func GetVIPUserCount() int64 {
	return DB.Where("is_v_ip = ?", true).Find(&model.User{}).RowsAffected
}

// GetToDayUpload 获取上传变量数量
func GetToDayUpload(s, e string) int64 {
	var r []model.Record
	return DB.Where("created_at between ? and ?", s, e).Find(&r).RowsAffected
}
