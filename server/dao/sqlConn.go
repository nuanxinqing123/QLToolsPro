// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 20:58
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : sqlConn.go

package dao

import (
	"QLToolsPro/server/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"os"
	"time"
)

var DB *gorm.DB
var err error
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Init() {
	// 配置日志
	var newLogger logger.Interface
	if viper.GetString("app.mode") == "debug" {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				LogLevel:                  logger.Warn, // 日志级别
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				LogLevel:                  logger.Error, // 日志级别
				IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,         // 禁用彩色打印
			},
		)
	}

	// 连接MySQL
	DB, err = gorm.Open(sqlite.Open("config/app.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		zap.L().Error("SQLite 发生错误, 原因：" + err.Error())
		panic(err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&model.User{},
		&model.Env{},
		&model.Panel{},
		&model.JWTAdmin{},
		&model.OperationRecord{},
		&model.CDK{},
		&model.Recharge{},
		&model.Record{},
		&model.WebSettings{},
		&model.Task{})

	if err != nil {
		zap.L().Error("SQLite 自动迁移失败, 原因：" + err.Error())
		panic(err.Error())
	}

	return
}

// InitWebSettings 初始化数据表
func InitWebSettings() {
	// 检查JWT密钥表是否存在
	jwtKey := GetJWTKey()
	if jwtKey == "" || len(jwtKey) < 10 {
		// 生成密码并写入数据库
		b := make([]rune, 18)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = letters[r.Intn(62)]
		}
		zap.L().Debug("生成密钥：" + string(b))
		CreateJWTKey(string(b))
	}

	// 判断Settings是否是第一次创建
	settings, err := GetSettings()
	if err != nil {
		zap.L().Error("InitWebSettings 发生错误")
		panic(err.Error())
	}
	if len(settings) == 0 {
		zap.L().Debug("Init WebSettings")
		p := &[]model.WebSettings{
			{Key: "web_title", Value: "青龙Tools Pro"},
			{Key: "web_ico", Value: "https://6b7.xyz/img/favicon.ico"},
			{Key: "web_logo", Value: "https://6b7.xyz/img/logo.png"},
			{Key: "register", Value: "1"},
			{Key: "notice", Value: "公告栏"},
		}

		err = SaveSettings(p)
		if err != nil {
			zap.L().Error("InitWebSettings 发生错误")
			panic(err.Error())
		}
	}
}
