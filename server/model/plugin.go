// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 20:07
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : plugin.go

package model

// DeletePlugin 删除插件
type DeletePlugin struct {
	TypeData string `json:"type" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

// FileInfo 读取普通插件信息
type FileInfo struct {
	FileName   string `json:"file_name"`
	PluginName string `json:"plugin_name"`
}

// FileCronInfo 读取定时插件信息
type FileCronInfo struct {
	// 插件文件名
	FileName string `json:"file_name"`
	// 插件名称
	PluginName string `json:"plugin_name"`
	// Cron定时
	CronTime string `json:"cron_time"`
	// 所需变量
	NeedEnvName string `json:"need_env_name"`
}
