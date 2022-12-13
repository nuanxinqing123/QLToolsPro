// -*- coding: utf-8 -*-
// @Time    : 2022/8/20 10:52
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : container.go

package logic

import (
	"QLToolsPro/server/cron"
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
	"time"
)

// ContainerTransfer 容器：迁移
func ContainerTransfer(p *model.ContainerOperation) (res.ResCode, string) {
	// 根据ID查询服务器信息
	StartData := dao.GetPanelDataByID(p.Start)
	EndData := dao.GetPanelDataByID(p.End)

	// 检查白名单
	if StartData.PanelURL == "" {
		return res.CodeContainerError, "发起的容器未在白名单内"
	}
	if EndData.PanelURL == "" {
		return res.CodeContainerError, "接收的容器未在白名单内"
	}

	// 获取Start面板全部信息
	zap.L().Debug("容器迁移：获取Start面板全部信息")
	url := panel.StringHTTP(StartData.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(StartData.PanelParams)
	allData, _ := requests.Requests("GET", url, "", StartData.PanelToken)

	// 绑定数据
	var token model.PanelEnvAll
	err := json.Unmarshal(allData, &token)
	if err != nil {
		zap.L().Error("[容器：迁移]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 向B容器上传变量
	zap.L().Debug("容器迁移：向End容器上传变量")
	go EnvUpload(token, EndData, "迁移任务(上传)")

	// 获取A容器所有变量ID
	idGroup := `[`
	for i := 0; i < len(token.Data); i++ {
		if token.Data[i].OID == "" {
			idGroup += strconv.Itoa(token.Data[i].ID) + `,`
		} else {
			idGroup += token.Data[i].OID + `,`
		}
	}
	idGroup = idGroup[:len(idGroup)-1]
	idGroup = idGroup + `]`

	// 删除A容器变量
	zap.L().Debug("容器迁移：删除Start容器变量")
	zap.L().Debug(idGroup)
	PanelEnvDelete(idGroup, StartData, "迁移任务(删除)")

	return res.CodeSuccess, "任务完成"
}

// ContainerCopy 容器：复制
func ContainerCopy(p *model.ContainerOperation) (res.ResCode, string) {
	// 根据ID查询服务器信息
	StartData := dao.GetPanelDataByID(p.Start)
	EndData := dao.GetPanelDataByID(p.End)

	// 检查白名单
	if StartData.PanelURL == "" {
		return res.CodeContainerError, "发起的容器未在白名单内"
	}
	if EndData.PanelURL == "" {
		return res.CodeContainerError, "接收的容器未在白名单内"
	}

	// 获取Start面板全部信息
	zap.L().Debug("容器复制：获取Start面板全部信息")
	url := panel.StringHTTP(StartData.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(StartData.PanelParams)
	allData, _ := requests.Requests("GET", url, "", StartData.PanelToken)

	// 绑定数据
	var token model.PanelEnvAll
	err := json.Unmarshal(allData, &token)
	if err != nil {
		zap.L().Error("[容器：复制]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 向B容器上传变量
	zap.L().Debug("容器复制：向B容器上传变量")
	go EnvUpload(token, EndData, "复制任务(上传)")

	return res.CodeSuccess, "任务完成"
}

// ContainerBackup 容器：备份
func ContainerBackup(p *model.ContainerOperationOne) (res.ResCode, string) {
	// 根据ID查询服务器信息
	StartData := dao.GetPanelDataByID(p.Start)

	// 检查白名单
	if StartData.PanelURL == "" {
		return res.CodeContainerError, "发起的容器未在白名单内"
	}

	// 获取Start面板全部信息
	zap.L().Debug("容器备份：获取面板全部信息")
	url := panel.StringHTTP(StartData.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(StartData.PanelParams)
	allData, _ := requests.Requests("GET", url, "", StartData.PanelToken)

	// 绑定数据
	var token model.PanelEnvAll
	err := json.Unmarshal(allData, &token)
	if err != nil {
		zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 创建JSON文件
	_, err = os.Create("backup.json")
	if err != nil {
		// 记录错误
		zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "创建备份文件失败"
	}

	// 打开JSON文件
	f, err := os.Open("backup.json")
	if err != nil {
		// 记录错误
		zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "打开备份文件失败"
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			// 记录错误
			zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		}
	}(f)

	// 序列化数据
	b, err := json.Marshal(token.Data)
	if err != nil {
		// 记录错误
		zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 保存数据
	err = os.WriteFile("backup.json", b, 0777)
	if err != nil {
		// 记录错误
		zap.L().Error("[容器：备份]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "保存备份文件失败"
	}

	return res.CodeSuccess, "任务完成"
}

// ContainerRestore 容器：恢复
func ContainerRestore(id string) (res.ResCode, string) {
	// 根据ID查询服务器信息
	iId, err := strconv.Atoi(id)
	if err != nil {
		zap.L().Error("[容器：恢复]查询面板信息失败:" + err.Error())
		return res.CodeServerBusy, "查询面板信息失败"
	}
	StartData := dao.GetPanelDataByID(iId)

	// 检查白名单
	if StartData.PanelURL == "" {
		return res.CodeContainerError, "发起的容器未在白名单内"
	}

	// 读取本地数据
	var backup model.PanelEnvAll
	// 打开文件
	file, err := os.Open("./backup.json")
	if err != nil {
		// 打开文件时发生错误
		zap.L().Error("[容器：恢复]打开备份文件失败:" + err.Error())
		return res.CodeServerBusy, "打开备份文件失败"
	}
	// 延迟关闭
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(file)

	// 配置读取
	byteData, err2 := io.ReadAll(file)
	if err2 != nil {
		// 读取配置时发生错误
		zap.L().Error("[容器：恢复]读取备份文件失败:" + err.Error())
		return res.CodeServerBusy, "读取备份文件失败"
	}

	// 数据绑定
	err3 := json.Unmarshal(byteData, &backup.Data)
	if err3 != nil {
		// 数据绑定时发生错误
		zap.L().Error("[容器：恢复]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 上传数据
	go EnvUpload(backup, StartData, "恢复任务")
	// 删除本地数据
	go DelBackupJSON()
	return res.CodeSuccess, "任务完成"
}

// ContainerSynchronization config.sh 同步
func ContainerSynchronization(p *model.ContainerOperationMass) (res.ResCode, string) {
	// 根据ID查询服务器信息
	StartData := dao.GetPanelDataByID(p.Start)

	// 检查白名单
	if StartData.PanelURL == "" {
		return res.CodeContainerError, "发起的容器未在白名单内"
	}

	// 获取Start面板的Config.sh
	zap.L().Debug("config.sh 同步：获取Start面板config.sh")
	url := panel.StringHTTP(StartData.PanelURL) + "/open/configs/config.sh?t=" + strconv.Itoa(StartData.PanelParams)
	allData, _ := requests.Requests("GET", url, "", StartData.PanelToken)

	// 绑定数据
	var sh model.ConfigSH
	err := json.Unmarshal(allData, &sh)
	if err != nil {
		zap.L().Error("[config.sh 同步]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}
	cd, err := json.Marshal(sh.Data)
	if err != nil {
		zap.L().Error("[config.sh 同步]序列化数据失败:" + err.Error())
		return res.CodeServerBusy, "序列化数据失败"
	}

	// 同步其余容器
	for i := 0; i < len(p.End); i++ {
		EndData := dao.GetPanelDataByID(p.End[i])
		if EndData.PanelURL == "" {
			continue
		}
		b := ConfigUpload(EndData, fmt.Sprintf("%s", cd))
		if b != true {
			return res.CodeContainerError, "面板：" + EndData.PanelName + "，Config.sh同步失败"
		}
	}
	return res.CodeSuccess, "任务完成"
}

// ConfigUpload Config.sh上传
func ConfigUpload(p model.Panel, content string) bool {
	data := `{"content": ` + content + `,"name": "config.sh"}`
	zap.L().Debug("Config.sh上传：" + data)
	url := panel.StringHTTP(p.PanelURL) + "/open/configs/save?t=" + strconv.Itoa(p.PanelParams)
	allData, _ := requests.Requests("POST", url, data, p.PanelToken)
	// 绑定数据
	var sh model.ConfigSH
	err := json.Unmarshal(allData, &sh)
	if err != nil {
		zap.L().Error("[config.sh 同步]序列化数据失败:" + err.Error())
	}
	if sh.Code != 200 {
		return false
	}
	return true
}

// EnvUpload 变量上传
func EnvUpload(p model.PanelEnvAll, pd model.Panel, journal string) {
	var re int
	var panelData model.Panel
	panelData = pd

	for i := 0; i < len(p.Data); i++ {
		var data string
		var pRes model.PanelRes

		if re == 1 {
			panelData = dao.GetPanelIDData(int(pd.ID))
		}

		zap.L().Debug("URL地址：" + panelData.PanelURL)
		URL := panel.StringHTTP(panelData.PanelURL) + "/open/envs?t=" + strconv.Itoa(panelData.PanelParams)
		// 上传
		data = `[{"value": "` + p.Data[i].Value + `","name": "` + p.Data[i].Name + `","remarks": "` + p.Data[i].Remarks + `"}]`
		// 执行上传任务
		r, err := requests.Requests("POST", URL, data, panelData.PanelToken)
		if err != nil {
			// 记录错误
			zap.L().Error("[容器：变量：上传]请求发送失败:" + err.Error())
			dao.RecordingError(journal, err.Error())
		}
		// 序列化内容
		err = json.Unmarshal(r, &pRes)
		if err != nil {
			zap.L().Error(err.Error())
		}

		if pRes.Code >= 401 && pRes.Code <= 500 {
			// 更新Token, 再次提交
			_, t := panel.GetPanelToken(panelData.PanelURL, panelData.PanelClientID, panelData.PanelClientSecret)
			pd.PanelToken = t.Data.Token
			pd.PanelParams = t.Data.Expiration
			re = 1
			i -= 1
		} else if pRes.Code == 400 {
			// 可能是重复上传，跳过
			continue
		} else if pRes.Code >= 500 {
			// 青龙请求错误, 再次提交
			i -= 1
		} else if pRes.Code == 200 {
			re = 0
		}

		// 限速
		time.Sleep(time.Second / 8)
	}
}

// PanelEnvDelete 变量删除
func PanelEnvDelete(p string, pd model.Panel, journal string) {
	URL := panel.StringHTTP(pd.PanelURL) + "/open/envs?t=" + strconv.Itoa(pd.PanelParams)
	// 执行删除任务
	_, err := requests.Requests("DELETE", URL, p, pd.PanelToken)
	if err != nil {
		// 记录错误
		dao.RecordingError(journal, err.Error())
	}
}

// DelBackupJSON 删除本地数据
func DelBackupJSON() {
	time.Sleep(time.Second * 5)
	err := os.Remove("backup.json")
	if err != nil {
		zap.L().Error(err.Error())
	}
}

// ContainerErrorContent 获取十条日志记录
func ContainerErrorContent() (res.ResCode, []model.OperationRecord) {
	// 查询记录
	return res.CodeSuccess, dao.ContainerErrorContent()
}

// ContainerCronBackup 定时备份面板变量
func ContainerCronBackup(p *model.CronBackUpEnv) (res.ResCode, string) {
	t := dao.NameGetCronTask("CronBackUpEnv")
	if t.Config == "" {
		// 未创建CronBackUpEnv任务，创建任务
		var c string
		for _, id := range p.PanelID {
			c += strconv.Itoa(id) + "&"
		}
		// 保存任务
		t.Name = "CronBackUpEnv"
		t.Cron = p.Cron
		t.Config = c
		t.State = p.State
		dao.SavaCronTask(t)

		// 暂停已启用任务
		cron.CStopTask()
		// 重启定时服务
		err := cron.CTask()
		if err != nil {
			return res.CodeContainerError, "重启定时任务服务"
		}
	} else {
		var c string
		for _, id := range p.PanelID {
			c += strconv.Itoa(id) + "&"
		}
		// 保存任务
		t.Cron = p.Cron
		t.Config = c
		t.State = p.State
		dao.SavaCronTask(t)

		// 暂停已启用任务
		cron.CStopTask()

		// 重启定时服务
		err := cron.CTask()
		if err != nil {
			return res.CodeContainerError, "重启定时任务服务"
		}
	}

	return res.CodeSuccess, "任务设置成功"
}
