// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 9:36
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : env.go

package dao

import (
	"QLToolsPro/server/model"
)

// GetDivisionEnvData 条件查询面板数据
func GetDivisionEnvData(page, q int) []model.Env {
	var env []model.Env
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Order("id desc").Limit(25).Offset(0).Find(&env)
	} else {
		// 获取第N页数据
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&env)
	}
	return env
}

// GetEnvDataPage 获取Env表总数据
func GetEnvDataPage() int64 {
	var c []model.Env
	result := DB.Find(&c)
	return result.RowsAffected
}

// EnvAdd 创建Env信息
func EnvAdd(data *model.EnvAdd) error {
	p := &model.Env{
		EnvName:       data.EnvName,
		EnvRemarks:    data.EnvRemarks,
		EnvQuantity:   data.EnvQuantity,
		EnvRegex:      data.EnvRegex,
		EnvMode:       data.EnvMode,
		EnvMerge:      data.EnvMerge,
		EnvUpdate:     data.EnvUpdate,
		EnvIsPlugin:   data.EnvIsPlugin,
		EnvPluginName: data.EnvPluginName,
		EnvIsCharging: data.EnvIsCharging,
	}
	return DB.Create(&p).Error
}

// EnvUpdate 更新Env信息
func EnvUpdate(data *model.EnvUpdate) error {
	var e model.Env
	DB.Where("id = ?", data.ID).First(&e)
	e.EnvName = data.EnvName
	e.EnvRemarks = data.EnvRemarks
	e.EnvQuantity = data.EnvQuantity
	e.EnvRegex = data.EnvRegex
	e.EnvMode = data.EnvMode
	e.EnvUpdate = data.EnvUpdate
	e.EnvIsPlugin = data.EnvIsPlugin
	e.EnvPluginName = data.EnvPluginName
	e.EnvIsCharging = data.EnvIsCharging
	e.EnvNeedIntegral = data.EnvNeedIntegral
	e.EnvTips = data.EnvTips
	e.EnvMerge = data.EnvMerge
	return DB.Save(&e).Error
}

// EnvDelete 删除Env信息
func EnvDelete(data *model.EnvDelete) error {
	return DB.Where("id IN ? ", data.ID).Delete(&model.Env{}).Error
}

// GetEnvAllData 获取面板全部变量数据
func GetEnvAllData() []model.Env {
	var e []model.Env
	DB.Find(&e)
	return e
}

// GetEnvNameCount 根据变量名获取配额
func GetEnvNameCount(name string) int {
	var env model.Env
	DB.Where("env_name = ?", name).First(&env)
	return env.EnvQuantity
}

// GetEnvNameData 根据变量名获取数据
func GetEnvNameData(name string) model.Env {
	var env model.Env
	DB.Where("env_name = ?", name).First(&env)
	return env
}

// GetIsChargingEnv 获取计费提交变量信息
func GetIsChargingEnv(c int) []model.Env {
	var env []model.Env
	DB.Where("env_is_charging = ?", c).Find(&env)
	return env
}
