// -*- coding: utf-8 -*-
// @Time    : 2022/8/26 10:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cdk.go

package dao

import (
	"QLToolsPro/server/model"
	"fmt"
	"go.uber.org/zap"
)

// CDKEYDivisionCDKData CD-KEY分页查询
func CDKEYDivisionCDKData(ctype string, page, q int) []model.CDK {
	var cdk []model.CDK
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Where("cd_key_type = ?", ctype).Order("id desc").Limit(25).Offset(0).Find(&cdk)
	} else {
		// 获取第N页数据
		DB.Where("cd_key_type = ?", ctype).Order("id desc").Limit(25).Offset((page - 1) * q).Find(&cdk)
	}
	return cdk
}

// GetCDKEYDataPage 获取CD-KEY表总数据
func GetCDKEYDataPage(ctype string) int64 {
	var c []model.CDK
	result := DB.Where("cd_key_type = ?", ctype).Find(&c)
	return result.RowsAffected
}

// CDKEYSearch CD-KEY数据筛选查询
func CDKEYSearch(ctype, s string) model.CDK {
	var c model.CDK
	DB.Where("cd_key_type = ? AND cd_key = ?", ctype, s).First(&c)
	return c
}

// CDKEYDataSearch CD-KEY数据查询
func CDKEYDataSearch(s string) model.CDK {
	var c model.CDK
	DB.Where("cd_key = ?", s).First(&c)
	return c
}

// CDKEYRemarksSearch CDK标识查询
func CDKEYRemarksSearch(s string) []model.CDK {
	var c []model.CDK
	DB.Where("cd_key_remarks LIKE ?", "%"+s+"%").Find(&c)
	return c
}

// CDKEYAdd 批量生成CDK
func CDKEYAdd(c *model.CDK) {
	var cdk model.CDK
	cdk.CdKey = c.CdKey
	cdk.CdKeyType = c.CdKeyType
	cdk.CdKeyIntegral = c.CdKeyIntegral
	cdk.CdKeyValidityPeriod = c.CdKeyValidityPeriod
	cdk.CdKeyRemarks = c.CdKeyRemarks
	cdk.CdKeyState = c.CdKeyState
	DB.Create(&cdk)
}

// CDKEYBatchUpdate 修改CD-KEY
func CDKEYBatchUpdate(p *model.UpdateCDK) error {
	zap.L().Debug("CdKeyState状态：" + fmt.Sprintf("%s", p.State))
	var cdk model.CDK
	DB.Where("id = ?", p.ID).First(&cdk)
	cdk.CdKeyRemarks = p.CdKeyRemarks
	cdk.CdKeyState = p.State
	return DB.Save(&cdk).Error
}

// CDKEYDelete 删除CD-KEY数据
func CDKEYDelete(p *model.DelCDK) error {
	var cdk model.CDK
	return DB.Where("id IN ?", p.IDList).Delete(&cdk).Error
}
