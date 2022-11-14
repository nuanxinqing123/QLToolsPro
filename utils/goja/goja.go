// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 21:34
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : goja.go

package goja

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	"QLToolsPro/utils/requests"
	"QLToolsPro/utils/wxpusher"
	"fmt"
	"github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonData struct {
	Bool bool   `json:"bool"`
	Env  string `json:"env"`
}

// RunPlugin 运行普通插件
func RunPlugin(filename, env string) (bool, string) {
	// 获取运行的绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}

	if !strings.Contains(filename, ".js") {
		// 不是JS文件
		return false, "输入文件类型错误"
	}

	// JS文件完整路径
	JSFilePath := ExecPath + "/plugin/ordinary/" + filename
	// 读取文件内容
	f, err := os.Open(JSFilePath)
	if err != nil {
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}
	data, err := io.ReadAll(f)
	if err != nil {
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}
	template := string(data)

	// 创建JS虚拟机
	vm := goja.New()
	// 注册JS方法
	vm.Set("request", request)
	vm.Set("console", console)
	vm.Set("refind", refind)
	vm.Set("replace", replace)
	vm.Set("mass_send_message", wxpusher.PluginSendMessageMass)
	vm.Set("sleep", sleep)
	_, err = vm.RunString(template)
	if err != nil {
		// JS代码有问题
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}
	var mainJs func(string) interface{}
	err = vm.ExportTo(vm.Get("main"), &mainJs)
	if err != nil {
		// JS函数映射到 Go函数失败
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}

	var j jsonData
	jd := mainJs(env)
	marshal, err := json.Marshal(jd)
	if err != nil {
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}
	err = json.Unmarshal(marshal, &j)
	if err != nil {
		zap.L().Error("[普通插件]：" + err.Error())
		return false, ""
	}

	zap.L().Debug(fmt.Sprintf("%v", j.Bool))
	zap.L().Debug(j.Env)

	return j.Bool, j.Env
}

// RunCronPlugin 运行定时插件
func RunCronPlugin(filename, envName string) {
	zap.L().Debug("定时插件名：" + filename)
	zap.L().Debug("定时插件所需变量：" + envName)

	// 查询所有面板此envName的变量
	env := GetPanelAppointEnvName(envName)

	// 获取运行的绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("[定时插件]：" + err.Error())
		return
	}

	if !strings.Contains(filename, ".js") {
		// 不是JS文件
		zap.L().Error("[定时插件]：输入文件类型错误")
		return
	}

	// JS文件完整路径
	JSFilePath := ExecPath + "/plugin/cron/" + filename
	// 读取文件内容
	f, err := os.Open(JSFilePath)
	if err != nil {
		zap.L().Error("[定时插件]：" + err.Error())
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		zap.L().Error("[定时插件]：" + err.Error())
		return
	}
	template := string(data)

	// 创建JS虚拟机
	vm := goja.New()
	// 注册JS方法
	vm.Set("request", request)
	vm.Set("console", console)
	vm.Set("refind", refind)
	vm.Set("replace", replace)
	vm.Set("mass_send_message", wxpusher.PluginSendMessageMass)
	vm.Set("sleep", sleep)
	_, err = vm.RunString(template)
	if err != nil {
		// JS代码有问题
		zap.L().Error("[定时插件]：" + err.Error())
		return
	}
	var mainJs func(interface{})
	err = vm.ExportTo(vm.Get("main"), &mainJs)
	if err != nil {
		// JS函数映射到 Go函数失败
		zap.L().Error("[定时插件]：" + err.Error())
		return
	}

	// 执行插件
	zap.L().Debug("执行插件")
	mainJs(env)
}

// GetPanelAppointEnvName 查询面板上指定变量名的变量
func GetPanelAppointEnvName(envName string) []model.EnvAllID {
	var eList []model.EnvAllID
	panelData := dao.GetPanelAllData()
	for _, m := range panelData {
		url := panel.StringHTTP(m.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(m.PanelParams)
		allData, _ := requests.Requests("GET", url, "", m.PanelToken)

		// 绑定数据
		var e model.EnvDataResult
		err := json.Unmarshal(allData, &e)
		if err != nil {
			zap.L().Error("[指定变量名查询]序列化数据失败:" + err.Error())
		}
		for _, eD := range e.Data {
			if eD.Name == envName && eD.Status == 0 {
				eList = append(eList, eD)
			}
		}
	}
	return eList
}
