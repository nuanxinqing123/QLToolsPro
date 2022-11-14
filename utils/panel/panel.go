// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 19:19
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : panel.go

package panel

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/gcache"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"QLToolsPro/utils/timeTools"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"strconv"
	time2 "time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StringHTTP 处理URL地址结尾的斜杠
func StringHTTP(url string) string {
	s := []byte(url)
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	url = string(s)
	return url
}

// GetPanelToken 获取面板Token(有效期：30天)
func GetPanelToken(url, id, secret string) (res.ResCode, model.Token) {
	var token model.Token

	URL := StringHTTP(url) + fmt.Sprintf("/open/auth/token?client_id=%s&client_secret=%s", id, secret)

	// 请求Token
	strData, err := requests.Requests("GET", URL, "", "")
	if err != nil {
		return res.CodeServerBusy, token
	}

	// 序列化内容
	err = json.Unmarshal(strData, &token)
	if err != nil {
		zap.L().Error(err.Error())
		return res.CodeServerBusy, token
	}

	// 更新数据库储存Token
	go dao.SaveToken(url, token.Data.Token, token.Data.Expiration)

	return res.CodeSuccess, token
}

// TestGetPanelToken 测试面板连接
func TestGetPanelToken(url, id, secret string) (res.ResCode, model.Token) {
	var token model.Token

	URL := fmt.Sprintf("/open/auth/token?client_id=%s&client_secret=%s", id, secret)
	nUrl := StringHTTP(url)

	// 请求Token
	strData, err := requests.Requests("GET", nUrl+URL, "", "")
	if err != nil {
		return res.CodeServerBusy, token
	}

	// 序列化内容
	err = json.Unmarshal(strData, &token)
	if err != nil {
		zap.L().Error(err.Error())
		return res.CodeServerBusy, token
	}

	return res.CodeSuccess, token
}

// UpdateOnlineServerCache 更新在线服务缓存
func UpdateOnlineServerCache() {
	// 实时计算缓存
	zap.L().Debug("实时计算在线服务")
	// 实时计算
	var os model.OnlineServer
	var OSData []byte
	// 获取启用服务器信息
	serverData := dao.GetPanelStartAllData()
	// 序列化数据
	data, err := json.Marshal(serverData)
	if err != nil {
		zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
		return
	}
	err = json.Unmarshal(data, &os.Online)
	if err != nil {
		zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
		return
	}
	//	统计数据
	for i := 0; i < len(os.Online); i++ {
		// 获取单节点容器数据
		thisEnvData := dao.GetPanelBindingEnv(int(serverData[i].ID))
		if len(thisEnvData) != 0 {
			//	序列化数据
			eData, err := json.Marshal(thisEnvData)
			if err != nil {
				zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
				return
			}
			err = json.Unmarshal(eData, &os.Online[i].EnvData)
			if err != nil {
				zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
				return
			}
		}

		// 获取面板已存在变量数量
		url := StringHTTP(serverData[i].PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(serverData[i].PanelParams)
		allData, err := requests.Requests("GET", url, "", serverData[i].PanelToken)
		if err != nil {
			zap.L().Error("面板：" + serverData[i].PanelName + ", 已无法连接，请管理员尽快处理")
			os.Online[i].PanelName = "服务器已失去连接（禁止提交）"
			os.Online[i].ID = -100
			continue
		}

		//	序列化数据
		var token model.EnvDataResult
		err = json.Unmarshal(allData, &token)
		if err != nil {
			zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
			return
		}

		// 判断返回状态
		if token.Code > 200 && token.Code < 500 {
			// 尝试获取授权
			go GetPanelToken(serverData[i].PanelURL, serverData[i].PanelClientID, serverData[i].PanelClientSecret)

			// 未授权或Token失效
			i += 1
			continue
		} else if token.Code >= 500 {
			os.Online[i].PanelName = "服务器数据处理失败"
			os.Online[i].ID = -100
		}

		// 计算变量剩余限额
		for x := 0; x < len(os.Online[i].EnvData); x++ {
			if thisEnvData[x].EnvMode == 1 || thisEnvData[x].EnvMode == 2 {
				// 新建模式
				_, _, os.Online[i].EnvData[x].EnvQuantity = CalculateQuantity(os.Online[i].ID, thisEnvData[x].EnvMode, os.Online[i].EnvData[x].EnvName)
			}
		}

		// 序列化数据
		OSData, err = json.Marshal(os)
		if err != nil {
			zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
			return
		}
	}
	//	缓存数据
	gcache.DeleteCache("online_server")
	gcache.TimingCache("online_server", string(OSData), time2.Hour*24)
}

// CalculateQuantity 计算变量剩余位置
func CalculateQuantity(id, mode int, name string) (res.ResCode, model.EnvDataResult, int) {
	var token model.EnvDataResult
	// 获取变量数据
	count := dao.GetEnvNameCount(name)

	// 获取容器信息
	sData := dao.GetPanelDataByID(id)

	// 获取面板已存在变量数量
	url := StringHTTP(sData.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(sData.PanelParams)
	allData, err := requests.Requests("GET", url, "", sData.PanelToken)
	if err != nil {
		return res.CodeServerBusy, token, 0
	}

	err = json.Unmarshal(allData, &token)
	if err != nil {
		zap.L().Error("[计算变量剩余位置]失败，原因：" + err.Error())
		return res.CodeServerBusy, token, 0
	}

	// 计算变量剩余限额
	c := count
	if mode == 1 || mode == 2 {
		// 新建模式
		for i := 0; i < len(token.Data); i++ {
			if token.Data[i].Name == name {
				c--
			}
		}
	}
	zap.L().Debug("计算结果：" + strconv.Itoa(c))
	return res.CodeSuccess, token, c
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

	panel := dao.GetPanelStartAllData()
	for i := 0; i < len(panel); i++ {
		// 检查Token是否已过期
		zap.L().Debug("检查过期CDK任务：检查Token是否已过期")
		var EnvRes model.EnvDataResult
		url := StringHTTP(panel[i].PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(panel[i].PanelParams)
		allData, err := requests.Requests("GET", url, "", panel[i].PanelToken)
		if err != nil {
			i -= 1
			count += 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time2.Now().Unix()))
				break
			}
			continue
		}
		if err = json.Unmarshal(allData, &EnvRes); err != nil {
			i -= 1
			count += 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time2.Now().Unix()))
				break
			}
			continue
		}
		if EnvRes.Code >= 400 && EnvRes.Code < 500 {
			// 尝试更新Token
			go GetPanelToken(panel[i].PanelURL, panel[i].PanelClientID, panel[i].PanelClientSecret)
			time2.Sleep(time2.Second)
			count += 1
			i -= 1
			if count == 5 {
				zap.L().Error("[禁用过期会员变量]任务失败，时间：" + timeTools.SwitchTimeStampToDataYearNowTome(time2.Now().Unix()))
				break
			}
			panel = dao.GetPanelStartAllData()
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
			if time2.Now().Unix() > user.ActivationTime.Unix() {
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
							if panel[i].PanelVersion {
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
				if panel[i].PanelVersion {
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
			idDateUrl := StringHTTP(panel[i].PanelURL) + "/open/envs/disable?t=" + strconv.Itoa(panel[i].PanelParams)
			_, _ = requests.Requests("PUT", idDateUrl, DisableID, panel[i].PanelToken)
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
