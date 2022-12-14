// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:17
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : env.go

package model

import "gorm.io/gorm"

// Env 变量表
type Env struct {
	gorm.Model
	EnvName         string // 环境变量名称
	EnvRemarks      string // 环境变量名称备注
	EnvQuantity     int    // 环境变量数量上限
	EnvRegex        string // 环境变量匹配正则
	EnvMode         int    // 环境变量模式[1：新建模式、2：更新模式、3：合并模式]
	EnvMerge        string // 合并模式分隔符
	EnvUpdate       string // 环境变量更新匹配正则（更新模式）
	EnvIsPlugin     bool   // 环境变量是否使用插件
	EnvPluginName   string // 绑定的插件名称
	EnvIsCharging   int    // 环境变量是否计费（0：不计费、1：VIP提交、2：积分提交）
	EnvNeedIntegral int    // 环境变量提交积分扣额
	EnvTips         string // 环境变量提示内容
}

// EnvPageData 面板分页数据
type EnvPageData struct {
	Page     int64 `json:"page"`
	PageData []Env `json:"pageData"`
}

// EnvSimpleData 变量简易数据
type EnvSimpleData struct {
	ID      int    `json:"id"`
	EnvName string `json:"name"`
}

// EnvAdd 变量新增
type EnvAdd struct {
	EnvName         string `json:"env_name"  binding:"required"`       // 环境变量名称
	EnvRemarks      string `json:"env_remarks"`                        // 环境变量名称备注
	EnvQuantity     int    `json:"env_quantity"  binding:"required"`   // 环境变量数量上限
	EnvRegex        string `json:"env_regex"`                          // 环境变量匹配正则
	EnvMode         int    `json:"env_mode"  binding:"required"`       // 环境变量模式[1：新建模式、2：更新模式、3：合并模式]
	EnvUpdate       string `json:"env_update"`                         // 环境变量更新匹配正则（更新模式）
	EnvIsPlugin     bool   `json:"env_is_plugin"`                      // 环境变量是否使用插件
	EnvPluginName   string `json:"env_plugin_name"`                    // 绑定的插件名称
	EnvIsCharging   int    `json:"env_is_charging" binding:"required"` // 环境变量是否计费（1：不计费、2：VIP提交、3：积分提交）
	EnvNeedIntegral int    `json:"env_need_integral"`                  // 环境变量积分提交需要多少积分
	EnvMerge        string `json:"env_merge"`                          // 环境变量分隔符（合并模式）
}

// EnvUpdate 面板更新
type EnvUpdate struct {
	ID              int    `json:"id" binding:"required"`              // 环境变量ID
	EnvName         string `json:"env_name" binding:"required"`        // 环境变量名称
	EnvRemarks      string `json:"env_remarks"`                        // 环境变量名称备注
	EnvQuantity     int    `json:"env_quantity" binding:"required"`    // 环境变量数量上限
	EnvRegex        string `json:"env_regex"`                          // 环境变量匹配正则
	EnvMode         int    `json:"env_mode" binding:"required"`        // 环境变量模式[1：新建模式、2：更新模式、3：合并模式]
	EnvUpdate       string `json:"env_update"`                         // 环境变量更新匹配正则（更新模式）
	EnvIsPlugin     bool   `json:"env_is_plugin"`                      // 环境变量是否使用插件
	EnvPluginName   string `json:"env_plugin_name"`                    // 绑定的插件名称
	EnvIsCharging   int    `json:"env_is_charging" binding:"required"` // 环境变量是否计费（0：不计费、1：VIP提交、2：积分提交）
	EnvNeedIntegral int    `json:"env_need_integral"`                  // 环境变量需要积分数量
	EnvTips         string `json:"env_tips"`                           // 环境变量提示内容
	EnvMerge        string `json:"env_merge"`                          // 环境变量分隔符（合并模式）
}

// EnvDelete 面板删除
type EnvDelete struct {
	ID []int `json:"id" binding:"required"` // 面板ID值
}

// EnvDataResult 面板变量数据
type EnvDataResult struct {
	Code int        `json:"code"`
	Data []EnvAllID `json:"data"`
}

// EnvAllID 全部环境变量
type EnvAllID struct {
	ID      int    `json:"id"`
	OldID   string `json:"_id"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks"`
	Status  int    `json:"status"`
}

// PanelEnvDelete 面板变量删除
type PanelEnvDelete struct {
	PanelName string `json:"panel_name" binding:"required"` // 面板名称
	ID        int    `json:"id"`                            // 变量ID值
	ID2       string `json:"_id"`                           // 变量ID值
}

// OnlineEnvUpload 上传变量
type OnlineEnvUpload struct {
	// 服务器ID
	ServerID int `json:"server_id" binding:"required"`
	// 变量名
	EnvName string `json:"env_name" binding:"required"`
	// 变量值
	EnvData string `json:"env_data" binding:"required"`
}
