// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:03
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cron.go

package cron

import (
	"QLToolsPro/utils/goja"
	res "QLToolsPro/utils/response"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

var err error
var c *cron.Cron

// Task 定时任务
func Task() error {
	// 刷新并启用启动任务
	c = cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600))) // 设置时区

	// 定时任务区

	// 定时更新面板Token（0 0 1/15 * *）
	_, err = c.AddFunc("0 0 1/15 * *", func() { // 每15天
		PanelTokenUpdate()
	})

	// 禁用会员到期变量
	_, err = c.AddFunc("0 0/4 * * *", func() { // 每四个小时执行一次
		DisableExpiredAccountsEnv()
	})

	// 运行定时插件
	resCode, fl := PluginCronData()
	if resCode != res.CodeServerBusy {
		if len(fl) != 0 {
			for _, i2 := range fl {
				zap.L().Debug("[运行定时插件]定时时间：" + i2.CronTime)
				_, err = c.AddFunc(i2.CronTime, func() {
					// 执行函数
					goja.RunCronPlugin(i2.FileName, i2.NeedEnvName)
				})
			}
		}
	} else {
		panic("定时插件列表读取失败")
	}

	// 定时任务结束

	if err != nil {
		return err
	}
	c.Start()
	return nil
}

// StopTask 暂停任务
func StopTask() {
	c.Stop()
}
