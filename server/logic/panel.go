// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:33
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : panel.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
	"strconv"
)

// PanelDivisionData 面板分页查询
func PanelDivisionData(page, quantity string) (res.ResCode, model.PanelPageData) {
	var data []model.Panel
	var pageData model.PanelPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.GetDivisionPanelData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.GetDivisionPanelData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.GetDivisionPanelData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetPanelDataPage()
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

// PanelDivisionDataSimple 面板简易数据
func PanelDivisionDataSimple() (res.ResCode, []model.PanelSimpleData) {
	var data []model.Panel
	var sd []model.PanelSimpleData

	// 获取所有面板
	data = dao.GetPanelAllData()
	for i := 0; i < len(data); i++ {
		var psd model.PanelSimpleData
		psd.ID = int(data[i].ID)
		psd.PanelName = data[i].PanelName
		sd = append(sd, psd)
	}

	return res.CodeSuccess, sd
}

// PanelAdd 面板新增
func PanelAdd(p *model.PanelAdd) (res.ResCode, string) {
	// 保存进数据库
	if err := dao.PanelAdd(p); err != nil {
		zap.L().Error("Error insert database, err:", zap.Error(err))
		return res.CodePanelError, "创建面板信息失败"
	}

	return res.CodeSuccess, ""
}

// PanelUpdate 更新面板信息
func PanelUpdate(p *model.PanelUpdate) (res.ResCode, string) {
	// 更新数据库
	if err := dao.PanelUpdate(p); err != nil {
		zap.L().Error("Error update database, err:", zap.Error(err))
		return res.CodePanelError, "更新面板信息失败"
	}
	return res.CodeSuccess, ""
}

// PanelDelete 删除面板信息
func PanelDelete(p *model.PanelDelete) (res.ResCode, string) {
	// 删除面板信息
	if err := dao.PanelDelete(p); err != nil {
		zap.L().Error("Error delete database, err:", zap.Error(err))
		return res.CodePanelError, "更新面板信息失败"
	}
	return res.CodeSuccess, ""
}

// PanelBindingUpdate 修改面板绑定变量
func PanelBindingUpdate(p *model.PanelBindingUpdate) (res.ResCode, string) {
	if err := dao.UpdatePanelEnvData(p); err != nil {
		zap.L().Error("Error update database, err:", zap.Error(err))
		return res.CodePanelError, "更新面板信息失败"
	}
	return res.CodeSuccess, ""
}

// PanelUnbindUpdate 面板解除变量绑定
func PanelUnbindUpdate() res.ResCode {
	// 获取所有面板信息
	panelData := dao.GetPanelAllData()
	for i := 0; i < len(panelData); i++ {
		panelData[i].EnvBinding = ""
		// 保存数据
		dao.PanelUnbindUpdate(panelData[i])
	}
	return res.CodeSuccess
}

// PanelTokenUpdate 面板批量更新Token
func PanelTokenUpdate() res.ResCode {
	// 获取所有面板信息
	panelData := dao.GetPanelAllData()
	for i := 0; i < len(panelData); i++ {
		// 更新Token
		go panel.GetPanelToken(panelData[i].PanelURL, panelData[i].PanelClientID, panelData[i].PanelClientSecret)
	}
	return res.CodeSuccess
}
