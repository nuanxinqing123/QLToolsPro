// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 11:00
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : container.go

package model

import "gorm.io/gorm"

// OperationRecord 日志记录数据表
type OperationRecord struct {
	gorm.Model
	OccurrenceTime string // 发生时间
	Operation      string // 操作方式
	Journal        string // 记录日志
}

// ContainerOperation 容器操作
type ContainerOperation struct {
	Start int `json:"start" binding:"required"`
	End   int `json:"end" binding:"required"`
}

// ContainerOperationOne 单容器操作
type ContainerOperationOne struct {
	Start int `json:"start" binding:"required"`
}

// ContainerOperationMass 多容器操作
type ContainerOperationMass struct {
	Start int   `json:"start" binding:"required"`
	End   []int `json:"end" binding:"required"`
}

// ConfigSH Config.sh Res
type ConfigSH struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

// ConfigSHReq Config.sh Req
type ConfigSHReq struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PanelEnvAll 面板全部变量数据
type PanelEnvAll struct {
	Code int      `json:"code"`
	Data []EnvAll `json:"data"`
}

// EnvAll 备份&恢复&返回Token
type EnvAll struct {
	ID      int    `json:"id"`
	OID     string `json:"_id"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks"`
}

// PanelAndEnvAll 携带Panel名称的Env
type PanelAndEnvAll struct {
	PanelName string   `json:"panel_name"`
	PanelEnv  []EnvAll `json:"panel_env"`
}

// UserEnvUpdate 用户修改变量
type UserEnvUpdate struct {
	PanelName string `json:"panel_name" binding:"required"`
	ID        int    `json:"id"`
	ID2       string `json:"_id"`
	Name      string `json:"name" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Remarks   string `json:"remarks" binding:"required"`
}
