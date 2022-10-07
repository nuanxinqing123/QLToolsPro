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
	"strconv"
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
	id.VIPUserCount = int(dao.GetVIPUserCount())

	// 获取今天的时间
	tt := 24
	e, _ := time.ParseDuration("-" + strconv.Itoa(tt) + "h")
	e1 := time.Now().Add(e)
	e2 := e1.Format("2006-01-02")
	id.ToDayUploadCount = int(dao.GetToDayUpload(e2+" 00:00:00", e2+" 23:59:59"))

	return res.CodeSuccess, id
}
