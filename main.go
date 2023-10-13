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
	_ "QLToolsPro/utils/daemon"
	"QLToolsPro/utils/snowflake"
	"QLToolsPro/utils/validator"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	// 加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("Viper init failed, err:%v\n", err)
		return
	}

	// 判断注册插件文件夹
	if viper.GetString("app.mode") == "" {
		bol := IFPlugin()
		if bol != true {
			fmt.Println("自动创建插件文件夹失败, 请手动在程序根目录创建 /plugin/ordinary 和 /plugin/cron 文件夹")
			return
		}
	}

	// 注册开机自启
	RegisterAndStartAutomatically()

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

	// 启动定时插件服务
	if err := cron.Task(); err != nil {
		fmt.Printf("Cron init failed, err:%v\n", err)
		return
	}
	zap.L().Debug("Cron success init ...")

	// 启动定时任务服务
	if err := cron.CTask(); err != nil {
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

	fmt.Println(" ")
	fmt.Println("  ____  _   _______          _     _____           \n / __ \\| | |__   __|        | |   |  __ \\          \n| |  | | |    | | ___   ___ | |___| |__) | __ ___  \n| |  | | |    | |/ _ \\ / _ \\| / __|  ___/ '__/ _ \\ \n| |__| | |____| | (_) | (_) | \\__ \\ |   | | | (_) |\n \\___\\_\\______|_|\\___/ \\___/|_|___/_|   |_|  \\___/ ")
	fmt.Println("")
	if viper.GetString("app.mode") == "debug" {
		fmt.Println("运行模式：Debug模式")
	} else {
		fmt.Println("运行模式：Release模式")
	}
	fmt.Println("监听端口：" + strconv.Itoa(viper.GetInt("app.port")))
	fmt.Println(" ")
	zap.L().Info("监听端口：" + strconv.Itoa(viper.GetInt("app.port")))

	//if err := endless.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), r); err != nil {
	//	log.Fatalf("listen: %s\n", err)
	//}
	//log.Println("Server exiting")

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	// 启动
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	// 十秒内优雅关闭服务（将未处理完成的请求处理完再关闭服务），超过十秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Service timed out has been shut down：", zap.Error(err))
	}

	zap.L().Info("Service has been shut down")

	// Linux
	//zap.L().Debug("进程PID：" + strconv.Itoa(syscall.Getpid()))
	//flag.Parse()
	//
	//listener, err := reload.GetListener(srv.Addr)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//var s = reload.NewService(listener)
	//go func() {
	//	err = srv.Serve(listener)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()
	//
	//s.Start()
}

// IFPlugin 判断并自动创建插件文件夹
func IFPlugin() bool {
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("[自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
		return false
	}

	_, err = os.Stat(ExecPath + "/plugin")
	if err != nil {
		err = os.Mkdir("plugin", 0777)
		if err != nil {
			fmt.Printf("Create Config Dir Error: %s", err)
			return false
		}
		_, err2 := os.Stat(ExecPath + "/plugin")
		if err2 != nil {
			zap.L().Error(err.Error())
			return false
		}
	}

	_, err = os.Stat(ExecPath + "/plugin/ordinary")
	if err != nil {
		err = os.Mkdir("plugin/ordinary", 0777)
		if err != nil {
			fmt.Println("[1、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
		_, err2 := os.Stat(ExecPath + "/plugin/ordinary")
		if err2 != nil {
			fmt.Println("[2、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
	}

	_, err = os.Stat(ExecPath + "/plugin/cron")
	if err != nil {
		err = os.Mkdir("plugin/cron", 0777)
		if err != nil {
			fmt.Println("[1、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
		_, err2 := os.Stat(ExecPath + "/plugin/cron")
		if err2 != nil {
			fmt.Println("[2、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
	}

	_, err = os.Stat(ExecPath + "/plugin/front")
	if err != nil {
		err = os.Mkdir("plugin/front", 0777)
		if err != nil {
			fmt.Println("[1、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
		_, err2 := os.Stat(ExecPath + "/plugin/front")
		if err2 != nil {
			fmt.Println("[2、自动创建插件文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
	}

	_, err = os.Stat(ExecPath + "/config/backup")
	if err != nil {
		err = os.Mkdir("config/backup", 0777)
		if err != nil {
			fmt.Println("[1、自动创建备份文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
		_, err2 := os.Stat(ExecPath + "/config/backup")
		if err2 != nil {
			fmt.Println("[2、自动创建备份文件夹]失败，" + fmt.Sprintf("原因：%v", err.Error()))
			return false
		}
	}

	return true
}

var ProcessName = getProcessName()

var getProcessName = func() string {
	return regexp.MustCompile(`([^/\s]+)$`).FindStringSubmatch(os.Args[0])[1]
}

// RegisterAndStartAutomatically 注册开机自启
func RegisterAndStartAutomatically() {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		fmt.Println("Windows & Darwin 系统暂不支持创建开机自启")
		return
	}

	b, err1 := PathExists("/usr/lib/systemd/system/QLToolsPro.service")
	if err1 != nil {
		fmt.Println("检查开机自启配置文件失败，原因：" + err1.Error())
	}
	if !b {
		ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println("获取程序运行目录绝对路径错误：" + err.Error())
			return
		}

		service := `
[Unit]
Description=QLToolsPro Service
After=network.target nss-lookup.target
Wants=network.target

[Service]
User=root
Group=root
Type=simple
LimitAS=infinity
LimitRSS=infinity
LimitCORE=infinity
LimitNOFILE=999999
WorkingDirectory=` + ExecPath + "/" + `
ExecStart=` + ExecPath + "/" + ProcessName + `
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target`
		data, err := exec.Command("sh", "-c", "type systemctl").Output()
		if err != nil {
			fmt.Println(err)
			return
		}

		if !strings.Contains(string(data), "bin") {
			return
		}
		os.WriteFile("/usr/lib/systemd/system/QLToolsPro.service", []byte(service), 0o644)
		exec.Command("systemctl", "disable", "QLToolsPro.service").Output()
		exec.Command("systemctl", "enable", "QLToolsPro.service").Output()
	}
}

/*
[Unit]
Description=QLToolsPro Service
After=network.target nss-lookup.target
Wants=network.target

[Service]
User=root
Group=root
Type=simple
LimitAS=infinity
LimitRSS=infinity
LimitCORE=infinity
LimitNOFILE=999999
WorkingDirectory=/root/ql/
ExecStart=` + ExecPath + "/" + ProcessName + `
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
*/

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	// IsNotExist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
