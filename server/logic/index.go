// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 22:23
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : index.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
	"time"
)

// AdminIndexData 管理员首页数据
func AdminIndexData() (res.ResCode, model.IndexData) {
	/*
		1、当前生效的面板数量
		2、当前生效的变量数量
		3、网站注册用户数量
		4、网站会员用户数量
		5、今日消费积分数量
		6、今日上传变量数量
	*/
	var id model.IndexData
	id.PanelCount = int(dao.GetPanelCount())
	id.EnvCount = int(dao.GetEnvCount())
	id.UserCount = int(dao.GetRegisterUserCount())

	// 获取今天的时间
	e := time.Now().Format("2006-01-02")
	zap.L().Debug("[管理员首页数据]开始时间" + e + " 00:00:00 & 结束时间：" + e + " 23:59:59")
	id.VIPUserCount = int(dao.GetVIPUserCount(e))
	id.ToDayIntegral = dao.GetIntegralToDayUpload(e+" 00:00:00", e+" 23:59:59")
	id.ToDayUploadCount = int(dao.GetEnvToDayUpload(e+" 00:00:00", e+" 23:59:59"))

	return res.CodeSuccess, id
}
