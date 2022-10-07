// -*- coding: utf-8 -*-
// @Time    : 2022/8/27 11:22
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : record.go

package model

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	RecordUserID     string // 上传记录用户ID
	RecordEnvName    string // 上传记录变量名
	RecordType       int    // 消费方式
	ConsumptionCount int    // 消费积分
}

// RecordPageData Record分页数据
type RecordPageData struct {
	Page     int64    `json:"page"`
	PageData []Record `json:"pageData"`
}
