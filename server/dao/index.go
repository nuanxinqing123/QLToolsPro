// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 23:17
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : index.go

package dao

import "QLToolsPro/server/model"

// GetPanelCount 获取面板存在数量
func GetPanelCount() int64 {
	return DB.Find(&[]model.Panel{}).RowsAffected
}

// GetEnvCount 获取变量存在数量
func GetEnvCount() int64 {
	return DB.Find(&[]model.Env{}).RowsAffected
}

// GetRegisterUserCount 获取网站注册用户数量
func GetRegisterUserCount() int64 {
	return DB.Find(&[]model.User{}).RowsAffected
}

// GetVIPUserCount 获取网站会员用户数量
func GetVIPUserCount(t string) int64 {
	return DB.Where("activation_time > ?", t).Find(&[]model.User{}).RowsAffected
}

// GetIntegralToDayUpload 获取今日消费积分
func GetIntegralToDayUpload(s, e string) int {
	var r []model.Record
	DB.Where("created_at between ? and ?", s, e).Find(&r)
	i := 0
	for _, re := range r {
		i += re.ConsumptionCount
	}
	return i
}

// GetEnvToDayUpload 获取上传变量数量
func GetEnvToDayUpload(s, e string) int64 {
	return DB.Where("created_at between ? and ?", s, e).Find(&[]model.Record{}).RowsAffected
}
