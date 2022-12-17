package model

import "gorm.io/gorm"

// Task 任务字段表
type Task struct {
	gorm.Model
	Name   string `json:"name,omitempty"`
	Cron   string `json:"cron,omitempty"`
	Config string `json:"config,omitempty"`
	State  bool   `json:"state"`
}

// TaskPageData Task分页数据
type TaskPageData struct {
	Page     int64  `json:"page"`
	PageData []Task `json:"pageData"`
}

// CronBackUpEnv 定时备份变量
type CronBackUpEnv struct {
	// 1 & 2 & 3 ...
	PanelID []int  `json:"panel_id" binding:"required"`
	Cron    string `json:"cron" binding:"required"`
	State   bool   `json:"state"`
}
