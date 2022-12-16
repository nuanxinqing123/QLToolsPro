// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:33
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : panel.go

package model

import "gorm.io/gorm"

// Panel 面板数据表
type Panel struct {
	gorm.Model
	PanelName         string // 面板名称
	PanelURL          string // 面板连接地址
	PanelClientID     string // 面板Client_ID
	PanelClientSecret string // 面板Client_Secret
	PanelEnable       bool   // 是否启用面板
	PanelVersion      bool   // 面板版本（false：新版本、true：旧版本）
	PanelToken        string // 面板Token
	PanelParams       int    // 面板Params
	EnvBinding        string // 面板绑定变量
}

// PanelPageData 面板分页数据
type PanelPageData struct {
	Page     int64   `json:"page"`
	PageData []Panel `json:"pageData"`
}

// PanelSimpleData 面板简易数据
type PanelSimpleData struct {
	ID        int    `json:"id"`
	PanelName string `json:"name"`
}

// PanelAdd 面板新增
type PanelAdd struct {
	Name         string `json:"name" binding:"required"`   // 面板名称
	URL          string `json:"url" binding:"required"`    // 面板连接地址
	ID           string `json:"id" binding:"required"`     // 面板Client_ID
	Secret       string `json:"secret" binding:"required"` // 面板Client_Secret
	Enable       bool   `json:"panel_enable"`              // 面板是否启用
	PanelVersion bool   `json:"panel_version"`             // 面板版本（false：新版本、true：旧版本）
}

// PanelUpdate 面板更新
type PanelUpdate struct {
	ID                int    `json:"id" binding:"required"`            // 面板ID值
	PanelName         string `json:"name" binding:"required"`          // 面板名称
	PanelURL          string `json:"url" binding:"required"`           // 面板连接地址
	PanelClientID     string `json:"client_id" binding:"required"`     // 面板Client_ID
	PanelClientSecret string `json:"client_secret" binding:"required"` // 面板Client_Secret
	PanelEnable       bool   `json:"panel_enable"`                     // 是否启用面板
	PanelVersion      bool   `json:"panel_version"`                    // 面板版本
}

// PanelDelete 面板删除
type PanelDelete struct {
	ID []int `json:"id" binding:"required"` // 面板ID值
}

// PanelBindingUpdate 面板修改绑定变量
type PanelBindingUpdate struct {
	ID         int      `json:"id" binding:"required"`          // 面板ID值
	EnvBinding []string `json:"env_binding" binding:"required"` // 变量值
}

// Token 面板Token数据
type Token struct {
	Code int `json:"code"`
	Data struct {
		Token      string `json:"token"`
		TokenType  string `json:"token_type"`
		Expiration int    `json:"expiration"`
	} `json:"data"`
	Message string
}

// PanelRes 面板Token数据
type PanelRes struct {
	Code    int `json:"code"`
	Message string
}

// OnlineServer 在线服务
type OnlineServer struct {
	// 可用服务器组
	Online []struct {
		// 容器ID
		ID int `json:"id"`
		// 容器名称
		PanelName string `json:"PanelName"`
		// 容器绑定变量
		EnvData []struct {
			// 变量名称
			EnvName string `json:"EnvName"`
			// 变量备注
			EnvRemarks string `json:"EnvRemarks"`
			// 变量剩余限额
			EnvQuantity int `json:"EnvQuantity"`
			// 变量提示内容
			EnvTips string `json:"EnvTips"`
		} `json:"env_data"`
	} `json:"online"`
}
