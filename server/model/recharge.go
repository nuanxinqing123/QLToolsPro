// -*- coding: utf-8 -*-
// @Time    : 2022/8/31 10:33
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : recharge.go

package model

import "gorm.io/gorm"

// Recharge 用户充值记录
type Recharge struct {
	gorm.Model
	RechargeUID  string // 充值用户ID
	RechargeCDK  string // 充值卡密
	RechargeType string // 充值类型
}

// RechargePage 用户充值记录分页数据
type RechargePage struct {
	Page     int64      `json:"page"`         // 总页数
	PageData []Recharge `json:"rechargeData"` // 分页查询数据
}

// UserRecharge 用户充值
type UserRecharge struct {
	RechargeCDK string `json:"recharge_cdk" binding:"required"`
}

// AdminRecharge 管理员充值
type AdminRecharge struct {
	UserID         string `json:"user_id" binding:"required"`
	RechargeType   int    `json:"recharge_type" binding:"required"`
	RechargeNumber int    `json:"recharge_number" binding:"required"`
}
