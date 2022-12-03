package model

import "gorm.io/gorm"

// Task 任务字段表
type Task struct {
	gorm.Model
	Name   string
	Cron   string
	Config string
	State  bool
}

// CronBackUpEnv 定时备份变量
type CronBackUpEnv struct {
	// 1 & 2 & 3 ...
	PanelID []int  `json:"panel_id" binding:"required"`
	Cron    string `json:"cron" binding:"required"`
	State   bool   `json:"state"`
}
