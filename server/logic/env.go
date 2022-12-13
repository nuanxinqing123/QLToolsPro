// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:16
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : env.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"strconv"
)

// 代替官方JSON库
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// EnvDivisionData 变量分页查询
func EnvDivisionData(page, quantity string) (res.ResCode, model.EnvPageData) {
	var data []model.Env
	var pageData model.EnvPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.GetDivisionEnvData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.GetDivisionEnvData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.GetDivisionEnvData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetEnvDataPage()
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

// EnvDivisionDataSimple 变量简易数据
func EnvDivisionDataSimple() (res.ResCode, []model.EnvSimpleData) {
	var data []model.Env
	var sd []model.EnvSimpleData

	// 获取所有面板
	data = dao.GetEnvAllData()
	for i := 0; i < len(data); i++ {
		var psd model.EnvSimpleData
		psd.ID = int(data[i].ID)
		psd.EnvName = data[i].EnvName
		sd = append(sd, psd)
	}

	return res.CodeSuccess, sd
}

// EnvAdd 面板新增
func EnvAdd(p *model.EnvAdd) (res.ResCode, string) {
	// 保存进数据库
	if err := dao.EnvAdd(p); err != nil {
		zap.L().Error("Error insert database, err:", zap.Error(err))
		return res.CodeEnvError, "创建变量信息失败"
	}

	return res.CodeSuccess, ""
}

// EnvUpdate 更新面板信息
func EnvUpdate(p *model.EnvUpdate) (res.ResCode, string) {
	// 更新数据库
	if err := dao.EnvUpdate(p); err != nil {
		zap.L().Error("Error update database, err:", zap.Error(err))
		return res.CodeEnvError, "更新变量信息失败"
	}
	return res.CodeSuccess, ""
}

// EnvDelete 删除面板信息
func EnvDelete(p *model.EnvDelete) (res.ResCode, string) {
	// 删除面板信息
	if err := dao.EnvDelete(p); err != nil {
		zap.L().Error("Error delete database, err:", zap.Error(err))
		return res.CodeEnvError, "更新变量信息失败"
	}
	return res.CodeSuccess, ""
}

// GetPanelAppointEnvName 查询面板上指定变量名的变量
//func GetPanelAppointEnvName(envName string) []model.EnvAllID {
//	var eList []model.EnvAllID
//	panelData := dao.GetPanelAllData()
//	for _, m := range panelData {
//		url := panel.StringHTTP(m.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(m.PanelParams)
//		allData, _ := requests.Requests("GET", url, "", m.PanelToken)
//
//		// 绑定数据
//		var e model.EnvAllID
//		err := json.Unmarshal(allData, &e)
//		if err != nil {
//			zap.L().Error("[指定变量名查询]序列化数据失败:" + err.Error())
//		}
//
//		if e.Name == envName {
//			eList = append(eList, e)
//		}
//	}
//
//	return eList
//}

// UserEnvDataSearch 用户变量管理:筛选查询
func UserEnvDataSearch(t, s string) (res.ResCode, []model.PanelAndEnvAll) {
	var ZpanelResult []model.PanelAndEnvAll
	if t == "panel" {
		var searchResult []interface{}
		var panelResult model.PanelAndEnvAll
		// 获取全部面板的变量
		iNum, err := strconv.Atoi(s)
		if err != nil {
			zap.L().Error("[用户变量管理:筛选查询]：查询内容类型错误")
			return res.CodeUserEnvError, ZpanelResult
		}
		p := dao.GetPanelIDData(iNum)
		// 检查请求面板是否存在
		if p.PanelName == "" {
			return res.CodeUserEnvError, ZpanelResult
		}
		// 获取面板上的所有变量
		url := panel.StringHTTP(p.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(p.PanelParams)
		allData, err := requests.Requests("GET", url, "", p.PanelToken)
		if err != nil {
			zap.L().Error("[用户变量管理:筛选查询]面板:" + p.PanelName + "。已无法连接，请管理员尽快处理")
			searchResult = append(searchResult, "查询面板:"+p.PanelName+"。可能已经失联，请管理员检查连通性")
			return res.CodeUserEnvError, ZpanelResult
		}
		var env model.PanelEnvAll
		err = json.Unmarshal(allData, &env)
		if err != nil {
			zap.L().Error("[用户变量管理:筛选查询]:" + err.Error())
			searchResult = append(searchResult, "数据序列化失败, 面板:"+p.PanelName)
			return res.CodeUserEnvError, ZpanelResult
		}

		// 判断返回状态
		if env.Code > 400 && env.Code < 500 {
			// 未授权或Token失效
			go panel.GetPanelToken(p.PanelURL, p.PanelClientID, p.PanelClientSecret)
			searchResult = append(searchResult, "查询面板:"+p.PanelName+"。可能已经失联，请管理员检查连通性")
			return res.CodeUserEnvError, ZpanelResult
		} else if env.Code > 500 {
			searchResult = append(searchResult, "查询面板:"+p.PanelName+"。发生错误，请检查青龙面板是否正常工作中")
			return res.CodeUserEnvError, ZpanelResult
		}
		// 添加符合规定的内容
		panelResult.PanelName = p.PanelName
		panelResult.PanelEnv = env.Data

		ZpanelResult = append(ZpanelResult, panelResult)
		return res.CodeSuccess, ZpanelResult
	} else {
		// 获取全部面板的变量
		for _, p := range dao.GetPanelStartAllData() {
			var searchResult []interface{}
			var panelResult model.PanelAndEnvAll

			// 获取面板上的所有变量
			url := panel.StringHTTP(p.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(p.PanelParams)
			allData, err := requests.Requests("GET", url, "", p.PanelToken)
			if err != nil {
				zap.L().Error("[用户变量管理:筛选查询]面板:" + p.PanelName + "。已无法连接，请管理员尽快处理")
				searchResult = append(searchResult, "查询面板:"+p.PanelName+"。可能已经失联，请管理员检查连通性")
				return res.CodeUserEnvError, ZpanelResult
			}
			var env model.PanelEnvAll
			err = json.Unmarshal(allData, &env)
			if err != nil {
				zap.L().Error("[用户变量管理:筛选查询]:" + err.Error())
				searchResult = append(searchResult, "数据序列化失败, 面板:"+p.PanelName)
				return res.CodeUserEnvError, ZpanelResult
			}

			// 判断返回状态
			if env.Code > 400 && env.Code < 500 {
				// 未授权或Token失效
				go panel.GetPanelToken(p.PanelURL, p.PanelClientID, p.PanelClientSecret)
				searchResult = append(searchResult, "查询面板:"+p.PanelName+"。可能已经失联，请管理员检查连通性")
				return res.CodeUserEnvError, ZpanelResult
			} else if env.Code > 500 {
				searchResult = append(searchResult, "查询面板:"+p.PanelName+"。发生错误，请检查青龙面板是否正常工作中")
				return res.CodeUserEnvError, ZpanelResult
			}

			// 添加符合规定的内容
			panelResult.PanelName = p.PanelName
			for _, e := range env.Data {
				if e.Remarks == s {
					panelResult.PanelEnv = append(panelResult.PanelEnv, e)
				}
			}
			ZpanelResult = append(ZpanelResult, panelResult)
		}
		return res.CodeSuccess, ZpanelResult
	}
}

// UserEnvDataUpdate 用户变量管理:修改用户变量
func UserEnvDataUpdate(p *model.UserEnvUpdate) (res.ResCode, string) {
	panelData := dao.GetPanelNameData(p.PanelName)
	if panelData.PanelURL == "" && panelData.PanelName == "" {
		return res.CodeUserEnvError, "修改的面板信息不存在"
	}

	url := panel.StringHTTP(panelData.PanelURL) + "/open/envs?t=" + strconv.Itoa(panelData.PanelParams)
	// 判断新旧面板版本
	data := ""
	if p.ID2 != "" {
		// 旧面板
		data = `{"_id": "` + p.ID2 + `", "value": "` + p.Value + `", "name": "` + p.Name + `", "remarks": "` + p.Remarks + `"}`
	} else {
		// 新面板
		data = `{"id": "` + strconv.Itoa(p.ID) + `", "value": "` + p.Value + `", "name": "` + p.Name + `", "remarks": "` + p.Remarks + `"}`
	}

	// 提交修改
	r, err := requests.Requests("PUT", url, data, panelData.PanelToken)
	if err != nil {
		zap.L().Error("[用户变量管理:修改用户变量]提交修改失败，原因：" + err.Error())
		return res.CodeServerBusy, "服务繁忙"
	}
	// 序列化内容
	var PanelRes model.PanelRes
	err = json.Unmarshal(r, &PanelRes)
	if err != nil {
		zap.L().Error(err.Error())
		return res.CodeServerBusy, "服务繁忙"
	}

	if PanelRes.Code >= 400 && PanelRes.Code < 500 {
		// 尝试更新Token
		go panel.GetPanelToken(panelData.PanelURL, panelData.PanelClientID, panelData.PanelClientSecret)
		return res.CodeUserEnvError, "发生一点小意外，请重新提交"
	} else if PanelRes.Code >= 500 {
		return res.CodeUserEnvError, "提交数据发生错误，错误原因：" + PanelRes.Message
	}

	return res.CodeSuccess, "更新成功"
}

// UserEnvDataDelete 用户变量管理:删除用户变量
func UserEnvDataDelete(p *model.PanelEnvDelete) (res.ResCode, string) {
	panelData := dao.GetPanelNameData(p.PanelName)
	if panelData.PanelURL == "" && panelData.PanelName == "" {
		return res.CodeUserEnvError, "面板信息不存在"
	}

	url := panel.StringHTTP(panelData.PanelURL) + "/open/envs?t=" + strconv.Itoa(panelData.PanelParams)
	// 判断新旧面板版本
	data := ""
	if p.ID2 != "" {
		// 旧面板
		data = `["` + p.ID2 + `"]`
	} else {
		// 新面板
		data = `[` + strconv.Itoa(p.ID) + `]`
	}

	// 提交修改
	r, err := requests.Requests("DELETE", url, data, panelData.PanelToken)
	if err != nil {
		zap.L().Error("[用户变量管理:删除用户变量]提交修改失败，原因：" + err.Error())
		return res.CodeServerBusy, "服务繁忙"
	}
	// 序列化内容
	var PanelRes model.PanelRes
	err = json.Unmarshal(r, &PanelRes)
	if err != nil {
		zap.L().Error("[用户变量管理:删除用户变量]序列化数据，原因：" + err.Error())
		return res.CodeServerBusy, "服务繁忙"
	}

	if PanelRes.Code >= 400 && PanelRes.Code < 500 {
		// 尝试更新Token
		go panel.GetPanelToken(panelData.PanelURL, panelData.PanelClientID, panelData.PanelClientSecret)
		return res.CodeUserEnvError, "发生一点小意外，请重新提交"
	} else if PanelRes.Code >= 500 {
		return res.CodeUserEnvError, "提交数据发生错误，错误原因：" + PanelRes.Message
	}

	return res.CodeSuccess, "删除成功"
}
