// -*- coding: utf-8 -*-
// @Time    : 2022/8/31 10:25
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : recharge.go

package dao

import (
	"QLToolsPro/server/model"
)

// UpdateUserIntegral 用户充值积分
func UpdateUserIntegral(p model.User) {
	DB.Save(&p)
}

// UpdateFalseCDK 禁用卡密
func UpdateFalseCDK(p model.CDK) {
	DB.Save(&p)
}

// InsertUserRechargeLog 记录用户充值记录
func InsertUserRechargeLog(p *model.Recharge) {
	DB.Create(&p)
}

// RechargeDivisionData 充值：记录查询
func RechargeDivisionData(page, q int) []model.Recharge {
	var recharge []model.Recharge
	if page == 1 {
		// 获取前20条数据
		DB.Order("id desc").Limit(25).Offset(0).Find(&recharge)
	} else {
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&recharge)
	}

	return recharge
}

// GetRechargeDataPage 获取Recharge表总数据
func GetRechargeDataPage() int64 {
	var c []model.Recharge
	result := DB.Find(&c)
	return result.RowsAffected
}

// RechargeSearch 数据查询
func RechargeSearch(so string, is bool) []model.Recharge {
	var r []model.Recharge
	if is {
		// CDK搜索
		DB.Where("recharge_cdk = ?", so).Find(&r)
	} else {
		// UserID搜索
		DB.Where("recharge_uid = ?", so).Find(&r)
	}

	return r
}

// RechargeUserDivisionData 用户充值：记录查询
func RechargeUserDivisionData(uid any, page, q int) []model.Recharge {
	var recharge []model.Recharge
	if page == 1 {
		// 获取前20条数据
		DB.Where("recharge_uid = ?", uid).Order("id desc").Limit(25).Offset(0).Find(&recharge)
	} else {
		DB.Where("recharge_uid = ?", uid).Order("id desc").Limit(25).Offset((page - 1) * q).Find(&recharge)
	}

	return recharge
}

// RechargeUserCount 用户充值：数据量
func RechargeUserCount(uid any) int64 {
	var c []model.Recharge
	return DB.Where("recharge_uid = ?", uid).Find(&c).RowsAffected
}
