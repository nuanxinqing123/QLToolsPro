// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 20:37
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : main.go

package main

import (
	"QLToolsPro/server"
	"QLToolsPro/server/cron"
	"QLToolsPro/server/dao"
	"QLToolsPro/server/gcache"
	"QLToolsPro/server/logger"
	"QLToolsPro/server/settings"
	"QLToolsPro/utils/license"
	"QLToolsPro/utils/snowflake"
	"QLToolsPro/utils/validator"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("Viper init failed, err:%v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Logger init failed, err:%v\n", err)
		return
	}
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(zap.L())
	zap.L().Debug("Logger success init ...")

	// 初始化数据库
	dao.Init()
	dao.InitWebSettings()
	zap.L().Debug("SQLite success init ...")

	// 初始化翻译器
	if err := validator.InitTrans("zh"); err != nil {
		fmt.Printf("Validator init failed, err:%v\n", err)
		return
	}
	zap.L().Debug("Validator success init ...")

	// 初始化雪花ID算法
	if err := snowflake.Init(); err != nil {
		fmt.Printf("Snowflake init failed, err:%v\n", err)
		return
	}
	zap.L().Debug("Snowflake success init ...")

	// 启动定时服务
	if err := cron.Task(); err != nil {
		fmt.Printf("Cron init failed, err:%v\n", err)
		return
	}
	zap.L().Debug("Cron success init ...")

	// 注册缓存
	gcache.InitCache()
	zap.L().Debug("Gcache success init ...")

	// 配置运行模式
	if viper.GetString("app.mode") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 注册路由
	r := server.Setup()

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	fmt.Println(" ")
	fmt.Println("  ____  _   _______          _     _____           \n / __ \\| | |__   __|        | |   |  __ \\          \n| |  | | |    | | ___   ___ | |___| |__) | __ ___  \n| |  | | |    | |/ _ \\ / _ \\| / __|  ___/ '__/ _ \\ \n| |__| | |____| | (_) | (_) | \\__ \\ |   | | | (_) |\n \\___\\_\\______|_|\\___/ \\___/|_|___/_|   |_|  \\___/ ")
	fmt.Println("")
	if viper.GetString("app.mode") == "debug" {
		fmt.Println("运行模式：Debug模式")
	} else {
		fmt.Println("运行模式：Release模式")
	}
	fmt.Println("监听端口：" + strconv.Itoa(viper.GetInt("app.port")))
	/* 检查授权 */
	b, m := license.LoginLicense()
	if !b {
		fmt.Println("您的许可状态：" + m)
		return
	} else {
		fmt.Println("您的许可状态：正常, 到期时间：" + m)
	}
	fmt.Println(" ")
	zap.L().Info("监听端口：" + strconv.Itoa(viper.GetInt("app.port")))

	// 启动
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listten: %s\n", err)
		}
	}()

	// 等待终端信号来优雅关闭服务器，为关闭服务器设置5秒超时
	quit := make(chan os.Signal, 1) // 创建一个接受信号的通道

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞此处，当接受到上述两种信号时，才继续往下执行
	zap.L().Info("Service ready to shut down")

	// 创建五秒超时的Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 五秒内优雅关闭服务（将未处理完成的请求处理完再关闭服务），超过十秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Service timed out has been shut down：", zap.Error(err))
	}

	zap.L().Info("Service has been shut down")
}
