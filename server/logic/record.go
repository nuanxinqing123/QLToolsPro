// -*- coding: utf-8 -*-
// @Time    : 2022/8/27 11:22
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : record.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
	"strconv"
)

// RecordDivisionData 上传记录分页查询
func RecordDivisionData(page, quantity string) (res.ResCode, model.RecordPageData) {
	var data []model.Record
	var pageData model.RecordPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.RecordDivisionData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.RecordDivisionData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.RecordDivisionData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetRecordDataPage()
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

// RecordSearch Record数据查询
func RecordSearch(s string) (res.ResCode, model.Record) {
	return res.CodeSuccess, dao.RecordSearch(s)
}

// RecordUserDivisionData 用户上传记录分页查询
func RecordUserDivisionData(uid any, page, quantity string) (res.ResCode, model.RecordPageData) {
	var data []model.Record
	var pageData model.RecordPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.RecordUserDivisionData(uid, 1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.RecordUserDivisionData(uid, 1, q)
		} else {
			// 查询指定页数的数据
			data = dao.RecordUserDivisionData(uid, intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetRecordUserDataPage(uid)
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
