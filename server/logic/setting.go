// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 21:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : setting.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
)

// GetSettings 获取所有配置信息
func GetSettings() (interface{}, res.ResCode) {
	data, err := dao.GetSettings()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, res.CodeServerBusy
	}

	return data, res.CodeSuccess
}

// SaveSettings 保存网站信息
func SaveSettings(p *[]model.WebSettings) res.ResCode {
	if err := dao.SaveSettings(p); err != nil {
		zap.L().Error(err.Error())
		return res.CodeServerBusy
	}

	return res.CodeSuccess
}
