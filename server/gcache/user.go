// -*- coding: utf-8 -*-
// @Time    : 2022/8/30 19:41
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : user.go

package gcache

import (
	"github.com/bluele/gcache"
	"go.uber.org/zap"
	"time"
)

var GC gcache.Cache

// InitCache 注册缓存插件
func InitCache() {
	GC = gcache.New(20).
		ARC().
		Build()
}

// TimingCache 定时缓存
func TimingCache(key, data string, t time.Duration) {
	err := GC.SetWithExpire(key, data, t)
	if err != nil {
		zap.L().Error("[定时缓存]失败， 原因" + err.Error())
	}
}

// GetCache 获取缓存内容
func GetCache(key string) (interface{}, error) {
	value, err := GC.Get(key)
	if err != nil {
		// 缓存不存在或缓存已失效
		return "", err
	}
	return value, nil
}

// DeleteCache 删除缓存内容
func DeleteCache(key string) {
	GC.Remove(key)
}
