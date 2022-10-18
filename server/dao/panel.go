// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 17:33
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : panel.go

package dao

import (
	"QLToolsPro/server/model"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// GetDivisionPanelData 条件查询面板数据
func GetDivisionPanelData(page, q int) []model.Panel {
	var cdk []model.Panel
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Order("id desc").Limit(25).Offset(0).Find(&cdk)
	} else {
		// 获取第N页数据
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&cdk)
	}
	return cdk
}

// GetPanelDataPage 获取Panel表总数据
func GetPanelDataPage() int64 {
	var c []model.Panel
	result := DB.Find(&c)
	return result.RowsAffected
}

// PanelAdd 创建新面板信息
func PanelAdd(data *model.PanelAdd) error {
	if data.Name == "" {
		data.Name = "未命名"
	}

	p := &model.Panel{
		PanelName:         data.Name,
		PanelURL:          data.URL,
		PanelClientID:     data.ID,
		PanelClientSecret: data.Secret,
		PanelEnable:       data.Enable,
		PanelVersion:      data.PanelVersion,
	}
	return DB.Create(&p).Error
}

// PanelUpdate 更新面板信息
func PanelUpdate(p *model.PanelUpdate) error {
	var pa model.Panel
	DB.Where("id = ? ", p.ID).First(&pa)
	pa.PanelName = p.PanelName
	pa.PanelURL = p.PanelURL
	pa.PanelClientID = p.PanelClientID
	pa.PanelClientSecret = p.PanelClientSecret
	pa.PanelEnable = p.PanelEnable
	pa.PanelVersion = p.PanelVersion
	return DB.Save(&pa).Error
}

// PanelDelete 删除面板信息
func PanelDelete(data *model.PanelDelete) error {
	return DB.Where("id IN ? ", data.ID).Delete(&model.Panel{}).Error
}

// GetPanelNameData 根据面板名称查询面板信息
func GetPanelNameData(s string) model.Panel {
	var p model.Panel
	DB.Where("panel_name = ?", s).First(&p)
	return p
}

// GetPanelIDData 根据面板ID查询面板信息
func GetPanelIDData(id int) model.Panel {
	var p model.Panel
	DB.Where("id = ?", id).First(&p)
	return p
}

// UpdatePanelEnvData 修改面板绑定变量
func UpdatePanelEnvData(data *model.PanelBindingUpdate) error {
	// 通过ID查询并更新数据
	return DB.Model(&model.Panel{}).Where("id = ? ", data.ID).Update("env_binding", strings.Join(data.EnvBinding, "@")).Error
}

// GetPanelAllData 获取面板All信息
func GetPanelAllData() []model.Panel {
	var p []model.Panel
	DB.Find(&p)
	return p
}

// GetPanelStartAllData 获取启用面板信息
func GetPanelStartAllData() []model.Panel {
	var p []model.Panel
	DB.Where("panel_enable = ?", true).Find(&p)
	return p
}

// PanelUnbindUpdate 面板解除变量绑定
func PanelUnbindUpdate(p model.Panel) {
	DB.Save(&p)
}

// SaveToken 储存Token
func SaveToken(url, token string, params int) {
	DB.Model(&model.Panel{}).Where("panel_url = ?", url).Updates(model.Panel{
		PanelToken:  token,
		PanelParams: params,
	})
}

// GetPanelBindingEnv 根据ID值获取面板绑定变量
func GetPanelBindingEnv(id int) []model.Env {
	// ID值查询服务
	var s model.Panel
	DB.Where("id = ?", id).First(&s)
	// 转换切片
	envBind := strings.Split(s.EnvBinding, "@")
	// 切片转换int类型
	var e []int
	zap.L().Debug("面板为：" + s.PanelName)
	for i := 0; i < len(envBind); i++ {
		zap.L().Debug("面板绑定变量数据为：" + envBind[i])
		if envBind[i] != "" {
			ee, err := strconv.Atoi(envBind[i])
			if err != nil {
				zap.L().Error(err.Error())
			}
			e = append(e, ee)
		}
	}

	var env []model.Env
	// 根据绑定值查询变量数据
	if len(e) != 0 {
		DB.Find(&env, e)
	}
	return env
}
