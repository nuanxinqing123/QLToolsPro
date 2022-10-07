// -*- coding: utf-8 -*-
// @Time    : 2022/8/29 23:18
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : index.go

package model

type IndexData struct {
	/*
		1、当前生效的面板数量
		2、当前生效的变量数量
		3、网站注册用户数量
		4、网站会员用户数量
		5、今日消费积分数量
		6、今日上传变量数量
	*/
	PanelCount       int `json:"panel_count"`
	EnvCount         int `json:"env_count"`
	UserCount        int `json:"user_count"`
	VIPUserCount     int `json:"vip_user_count"`
	ToDayIntegral    int `json:"to_day_integral"`
	ToDayUploadCount int `json:"to_day_upload_count"`
}
