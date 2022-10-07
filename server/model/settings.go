// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 10:58
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : setting.go

package model

// WebSettings 网站设置表
type WebSettings struct {
	Key   string `json:"key" gorm:"primaryKey" binding:"required"`
	Value string `json:"value"`
}
