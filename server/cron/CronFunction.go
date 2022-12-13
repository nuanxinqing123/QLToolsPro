package cron

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/panel"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"QLToolsPro/utils/timeTools"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ContainerBackup 定时备份容器
func ContainerBackup(c string) {
	// 获取程序路径
	ExecPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		zap.L().Error("[定时备份容器]失败:" + err.Error())
	}

	// 转换切片
	pid := strings.Split(c, "&")
	// 查询面板数据
	for _, p := range pid {
		pp, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		pd := dao.GetPanelDataByID(pp)

		zap.L().Debug("[定时备份容器]面板URL：" + pd.PanelURL)
		url := panel.StringHTTP(pd.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(pd.PanelParams)
		allData, _ := requests.Requests("GET", url, "", pd.PanelToken)

		// 绑定数据
		var token model.PanelEnvAll
		err = json.Unmarshal(allData, &token)
		if err != nil {
			zap.L().Error("[定时备份容器]序列化数据失败:" + err.Error())
		}

		// 创建JSON文件
		e := time.Now().Format("2006-01-02")
		path := ExecPath + "/config/backup/" + e + "_" + pd.PanelName + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".json"
		_, err = os.Create(path)
		if err != nil {
			// 记录错误
			zap.L().Error("[定时备份容器]创建数据文件失败:" + err.Error())
		}

		// 打开JSON文件
		f, err := os.Open(path)
		if err != nil {
			// 记录错误
			zap.L().Error("[定时备份容器]写入数据失败:" + err.Error())
		}
		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				// 记录错误
				zap.L().Error("[容器：备份]保存数据失败:" + err.Error())
			}
		}(f)

		// 序列化数据
		b, err := json.Marshal(token.Data)
		if err != nil {
			// 记录错误
			zap.L().Error("[定时备份容器]序列化数据失败:" + err.Error())
		}

		// 保存数据
		err = os.WriteFile(path, b, 0777)
		if err != nil {
			// 记录错误
			zap.L().Error("[定时备份容器]保存数据失败:" + err.Error())
		}
	}
}

// DisableExpiredAccountsEnv 禁用过期会员变量
func DisableExpiredAccountsEnv() {
	/*
		1、遍历所有面板，检查Token是否已过期（Token过期则更新，并重新检查）
		2、获取面板上所有状态（status）为启用的环境变量，并且属于会员提交的变量名
		3、获取所有用户信息，并去重
		4、校验用户会员是否到期，并记录
		5、禁用已过期会员变量
	*/
	var count int

	panels := dao.GetPanelStartAllData()
	for i := 0; i < len(panels); i++ {
		// 检查Token是否已过期
		zap.L().Debug("检查过期CDK任务：检查Token是否已过期")
		var EnvRes model.EnvDataResult
		url := StringHTTP(panels[i].PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(panels[i].PanelParams)
		allData, err := requests.Requests("GET", url, "", panels[i].PanelToken)
		if err != nil {
			i -= 1
			count += 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time.Now().Unix()))
				break
			}
			continue
		}
		if err = json.Unmarshal(allData, &EnvRes); err != nil {
			i -= 1
			count += 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time.Now().Unix()))
				break
			}
			continue
		}
		if EnvRes.Code >= 400 && EnvRes.Code < 500 {
			// 尝试更新Token
			go panel.GetPanelToken(panels[i].PanelURL, panels[i].PanelClientID, panels[i].PanelClientSecret)
			time.Sleep(time.Second)
			count += 1
			i -= 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time.Now().Unix()))
				break
			}
			panels = dao.GetPanelStartAllData()
			continue
		}

		// 获取会员计费变量信息
		var ea []string
		env := dao.GetIsChargingEnv(1)
		if len(env) == 0 {
			break
		}
		for _, ed := range EnvRes.Data {
			for _, en := range env {
				if ed.Status == 0 && en.EnvName == ed.Name {
					ea = append(ea, ed.Remarks)
				}
			}
		}

		// 去重
		ea = RemoveDuplicatesFromSlice(ea)

		// 检查用户会员是否到期
		var eaEnd []string
		for _, u := range ea {
			_, user := dao.GetUserNameData(u)
			if time.Now().Unix() > user.ActivationTime.Unix() {
				eaEnd = append(eaEnd, user.Username)
			}
		}

		// 禁用任务
		var id []string
		if len(eaEnd) != 0 {
			for _, en := range env {
				for _, ed := range EnvRes.Data {
					for _, eae := range eaEnd {
						if ed.Name == en.EnvName && ed.Remarks == eae {
							if panels[i].PanelVersion {
								// Old
								id = append(id, ed.OldID)
							} else {
								// New
								id = append(id, strconv.Itoa(ed.ID))
							}
						}
					}
				}
			}

			// 禁用
			DisableID := `[`
			for y := 0; y < len(id); y++ {
				if panels[i].PanelVersion {
					// Old
					if y == len(id)-1 {
						DisableID += "\"" + id[y] + "\""
					} else {
						DisableID += "\"" + id[y] + "\", "
					}
				} else {
					// New
					if y == len(id)-1 {
						DisableID += id[y]
					} else {
						DisableID += id[y] + ","
					}
				}
			}
			DisableID += `]`
			zap.L().Debug("检查会员过期变量任务：" + DisableID)
			idDateUrl := StringHTTP(panels[i].PanelURL) + "/open/envs/disable?t=" + strconv.Itoa(panels[i].PanelParams)
			_, _ = requests.Requests("PUT", idDateUrl, DisableID, panels[i].PanelToken)
		}

		zap.L().Info("[禁用过期会员变量]任务完成")
	}
}

// RemoveDuplicatesFromSlice 切片去重
func RemoveDuplicatesFromSlice(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// StringHTTP 处理URL地址结尾的斜杠
func StringHTTP(url string) string {
	s := []byte(url)
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	url = string(s)
	return url
}

// PanelTokenUpdate 面板批量更新Token
func PanelTokenUpdate() res.ResCode {
	// 获取所有面板信息
	panelData := dao.GetPanelAllData()
	for i := 0; i < len(panelData); i++ {
		// 更新Token
		go panel.GetPanelToken(panelData[i].PanelURL, panelData[i].PanelClientID, panelData[i].PanelClientSecret)
	}
	return res.CodeSuccess
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
	files, _ := os.ReadDir(PluginPath)

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
		v, _ := io.ReadAll(fd)
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
