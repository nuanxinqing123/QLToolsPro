package cron

import (
	"QLToolsPro/server/dao"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

var errTask error
var cTask *cron.Cron

// CTask 定时任务
func CTask() error {
	// 刷新并启用启动任务
	cTask = cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600))) // 设置时区

	// 获取任务列表数据
	t := dao.GetAllCronTask()

	// 定时任务区

	for _, task := range t {
		if task.State != false {
			switch task.Name {
			case "CronBackUpEnv":
				// 定时备份面板变量
				_, errTask = cTask.AddFunc(task.Cron, func() {
					ContainerBackup(task.Config)
				})
			default:
				zap.L().Info("暂无定时任务")
			}
		}
	}

	// 定时任务结束

	if errTask != nil {
		zap.L().Error("[系统定时任务]刷新失败，原因：" + errTask.Error())
		return errTask
	}
	cTask.Start()
	return nil
}

// CStopTask 暂停任务
func CStopTask() {
	cTask.Stop()
}
