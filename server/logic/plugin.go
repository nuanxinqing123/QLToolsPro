// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 19:52
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : plugin.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PluginOrdinaryData 普通插件查询
func PluginOrdinaryData() (res.ResCode, []model.FileInfo) {
	// 获取插件目录绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("获取插件目录绝对路径失败：" + err.Error())
		return res.CodeServerBusy, nil
	}

	// 普通插件目录
	PluginPath := ExecPath + "/plugin/ordinary/"

	// 读取目录
	var fl []model.FileInfo
	var fi model.FileInfo
	files, _ := ioutil.ReadDir(PluginPath)

	// 读取插件信息
	for _, f := range files {
		// 跳过不是JS的文件
		if !strings.Contains(f.Name(), ".js") {
			continue
		}

		zap.L().Debug("文件名：" + f.Name())

		// 读取插件名称
		fd, err2 := os.Open(PluginPath + f.Name())
		if err2 != nil {
			zap.L().Error(f.Name() + "：打开文件失败，原因：" + err2.Error())
		}
		defer func(fd *os.File) {
			err3 := fd.Close()
			if err3 != nil {
				zap.L().Error("文件名：" + f.Name() + "  关闭文件失败，原因：" + err3.Error())
			}
		}(fd)
		v, _ := ioutil.ReadAll(fd)
		data := string(v)
		PluginName := ""
		if regs := regexp.MustCompile(`\[name:(.+)]`).FindStringSubmatch(data); len(regs) != 0 {
			PluginName = strings.Trim(regs[1], " ")
		}

		fi.FileName = f.Name()
		fi.PluginName = PluginName
		fl = append(fl, fi)
	}
	return res.CodeSuccess, fl
}

// PluginCronData 定时插件查询
func PluginCronData() (res.ResCode, []model.FileCronInfo) {
	// 获取插件目录绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("获取插件目录绝对路径失败：" + err.Error())
		return res.CodeServerBusy, nil
	}

	// 定时插件目录
	PluginPath := ExecPath + "/plugin/cron/"

	// 读取目录
	var fl []model.FileCronInfo
	var fi model.FileCronInfo
	files, _ := ioutil.ReadDir(PluginPath)

	// 读取插件信息
	for _, f := range files {
		// 跳过不是JS的文件
		if !strings.Contains(f.Name(), ".js") {
			continue
		}

		zap.L().Debug("文件名：" + f.Name())

		// 读取插件名称
		fd, err2 := os.Open(PluginPath + f.Name())
		if err2 != nil {
			zap.L().Error(f.Name() + "：打开文件失败，原因：" + err2.Error())
		}
		defer func(fd *os.File) {
			err3 := fd.Close()
			if err3 != nil {
				zap.L().Error("文件名：" + f.Name() + "  关闭文件失败" + err3.Error())
			}
		}(fd)
		v, _ := ioutil.ReadAll(fd)
		data := string(v)
		PluginName := ""
		CronTime := ""
		NeedEnvName := ""
		if regs := regexp.MustCompile(`\[name:(.+)]`).FindStringSubmatch(data); len(regs) != 0 {
			PluginName = strings.Trim(regs[1], " ")
		}
		if cron := regexp.MustCompile(`\[cron:([^\[\]]+)]`).FindStringSubmatch(data); len(cron) != 0 {
			CronTime = strings.Trim(cron[1], " ")
		}
		if env := regexp.MustCompile(`\[env:([^\[\]]+)]`).FindStringSubmatch(data); len(env) != 0 {
			NeedEnvName = strings.Trim(env[1], " ")
		}

		fi.FileName = f.Name()
		fi.PluginName = PluginName
		fi.CronTime = CronTime
		fi.NeedEnvName = NeedEnvName
		fl = append(fl, fi)
	}
	return res.CodeSuccess, fl
}

// PluginDelete 删除插件
func PluginDelete(p *model.DeletePlugin) (res.ResCode, string) {
	// 检查插件是否还存在绑定
	env := dao.GetEnvAllData()
	ee := 0
	for _, e := range env {
		if e.EnvIsPlugin {
			if e.EnvPluginName == p.FileName {
				ee += 10
			}
		}
	}
	if ee != 0 {
		return res.CodePluginError, "欲删除插件任存在绑定，请解除绑定后再删除"
	}

	// 获取插件目录绝对路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("[删除插件:]" + err.Error())
		return res.CodePluginError, "删除插件失败"
	}

	// 删除文件
	FilePath := ""
	if p.TypeData == "ordinary" {
		FilePath = ExecPath + "/plugin/ordinary/" + p.FileName
	} else {
		FilePath = ExecPath + "/plugin/cron/" + p.FileName
	}
	err = os.Remove(FilePath)
	if err != nil {
		// 删除失败
		zap.L().Error("[删除插件:]" + err.Error())
		return res.CodePluginError, "删除插件失败"
	}

	return res.CodeSuccess, ""
}
