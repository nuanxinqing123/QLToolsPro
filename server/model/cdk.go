// -*- coding: utf-8 -*-
// @Time    : 2022/8/26 10:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cdk.go

package model

import (
	"gorm.io/gorm"
)

type CDK struct {
	gorm.Model
	CdKey               string // CD-KEY值
	CdKeyType           string // CD-KEY类型
	CdKeyIntegral       int    // CD-KEY积分
	CdKeyValidityPeriod int64  // CD-KEY有效期
	CdKeyRemarks        string // CD-KEY标识
	CdKeyState          bool   // CD-KEY状态（true：启用、false：禁用）
}

// CDKPageData CDK分页数据
type CDKPageData struct {
	Page     int64 `json:"page"`
	PageData []CDK `json:"pageData"`
}

// CreateCDK 生成CDK
// 有效期(1：天卡、7：周卡、14：两周卡、31：月卡、92：季度、183：半年、365：一年)
type CreateCDK struct {
	CdKeyCount          int    `json:"cd_key_count" binding:"required"` // CD-KEY生成数量
	CdKeyType           string `json:"cd_key_type" binding:"required"`  // CD-KEY类型
	CdKeyIntegral       int    `json:"cd_key_integral"`                 // CD-KEY积分
	CdKeyValidityPeriod int    `json:"cd_key_validity_period"`          // CD-KEY有效期
	CdKeyPrefix         string `json:"cd_key_prefix"`                   // CD-KEY前缀
	CdKeyRemarks        string `json:"cd_key_remarks"`                  // CD-KEY标识
}

// OutputCDK 导出CDK
type OutputCDK struct {
	/*
		wAll：所有未使用
		wMon：未使用月卡
		wQua：未使用季卡
		wHalf：未使用半年卡
		wYear：未使用年卡
		yAll：所有已使用
		yMon：已使用月卡
		yQua：已使用季卡
		yHalf：已使用半年卡
		yYear：已使用年卡
	*/

	CdKeyType string `json:"cdKeyType" binding:"required"` // CDK类型
}

// UpdateCDK 更新CDK数据
type UpdateCDK struct {
	ID           int    `json:"id" binding:"required"` // 操作CDK的ID
	State        bool   `json:"state"`                 // CDK状态
	CdKeyRemarks string `json:"cd_key_remarks"`        // CD-KEY标识
}

// BatchUpdateCDK 更新多条CDK数据
type BatchUpdateCDK struct {
	IDList       []int  `json:"idList" binding:"required"` // 操作CDK的ID
	State        bool   `json:"state"`                     // CDK状态
	CdKeyRemarks string `json:"cd_key_remarks"`            // CD-KEY标识
}

// DelCDK 删除CDK数据
type DelCDK struct {
	IDList []int `json:"id_list" binding:"required"` // CDK ID
}

// CheckCDK 检查CD-KEY
type CheckCDK struct {
	CDK string `json:"cdk" binding:"required"`
}
