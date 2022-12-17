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

// GetCronTask 获取定时任务
func GetCronTask(s string) model.Task {
	var t model.Task
	DB.Where("name = ?", s).First(&t)
	return t
}

// CronTaskDivisionData 分页查询
func CronTaskDivisionData(page, q int) []model.Task {
	var cdk []model.Task
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Order("id desc").Limit(25).Offset(0).Find(&cdk)
	} else {
		// 获取第N页数据
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&cdk)
	}
	return cdk
}

// CronTaskDivisionDataPage 查询总页数
func CronTaskDivisionDataPage() int64 {
	return DB.Find(&[]model.Task{}).RowsAffected
}
