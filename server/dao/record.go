// -*- coding: utf-8 -*-
// @Time    : 2022/8/27 11:22
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : record.go

package dao

import "QLToolsPro/server/model"

// RecordDivisionData 上传记录分页查询
func RecordDivisionData(page, q int) []model.Record {
	var record []model.Record
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Order("id desc").Limit(25).Offset(0).Find(&record)
	} else {
		// 获取第N页数据
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&record)
	}
	return record
}

// GetRecordDataPage 获取Record表总数据
func GetRecordDataPage() int64 {
	var c []model.Record
	result := DB.Find(&c)
	return result.RowsAffected
}

// RecordSearch 数据查询
func RecordSearch(s string) model.Record {
	var c model.Record
	DB.Where("record_user_id = ?", s).First(&c)
	return c
}

// RecordUserDivisionData 用户上传：记录查询
func RecordUserDivisionData(uid any, page, q int) []model.Record {
	var record []model.Record
	if page == 1 {
		// 获取前20条数据
		DB.Where("record_user_id = ?", uid).Order("id desc").Limit(25).Offset(0).Find(&record)
	} else {
		DB.Where("record_user_id = ?", uid).Order("id desc").Limit(25).Offset((page - 1) * q).Find(&record)
	}

	return record
}

// GetRecordUserDataPage 获取Record表用户数据
func GetRecordUserDataPage(uid any) int64 {
	var c []model.Record
	return DB.Where("record_user_id = ?", uid).Find(&c).RowsAffected
}

// CreateRecordInfo 添加上传记录
func CreateRecordInfo(uid any, envName string, upType, integral int) {
	var re model.Record
	re.RecordUserID = uid.(string)
	re.RecordEnvName = envName
	re.RecordType = upType
	re.ConsumptionCount = integral
	DB.Create(&re)
}
