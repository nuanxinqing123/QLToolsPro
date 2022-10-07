// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 11:01
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : container.go

package dao

import (
	"QLToolsPro/server/model"
	"QLToolsPro/utils/timeTools"
	"gorm.io/gorm"
	"time"
)

// RecordingError 记录错误
func RecordingError(journal, info string) {
	DB.Create(&model.OperationRecord{
		Model:          gorm.Model{},
		OccurrenceTime: timeTools.SwitchTimeStampToDataYearNowTome(time.Now().Unix()),
		Operation:      info,
		Journal:        journal,
	})
}

// GetPanelDataByID 根据ID值查询容器信息
func GetPanelDataByID(id int) model.Panel {
	var d model.Panel
	// 通过ID查询容器
	DB.Where("id = ?", id).First(&d)
	return d
}

// ContainerErrorContent 获取十条错误记录
func ContainerErrorContent() []model.OperationRecord {
	var i []model.OperationRecord
	DB.Order("id desc").Limit(10).Offset(0).Find(&i)
	return i
}
