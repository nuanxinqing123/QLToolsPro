package dao

import "QLToolsPro/server/model"

// NameGetCronTask 任务名获取定时配置文件
func NameGetCronTask(s string) model.Task {
	var t model.Task
	DB.Where("name = ?", s).First(&t)
	return t
}

// SavaCronTask 保存任务配置
func SavaCronTask(t model.Task) {
	DB.Save(&t)
}

// GetAllCronTask 获取全部定时任务
func GetAllCronTask() []model.Task {
	var t []model.Task
	DB.Find(&t)
	return t
}
