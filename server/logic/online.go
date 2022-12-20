// -*- coding: utf-8 -*-
// @Time    : 2022/9/3 12:00
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : online.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/gcache"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/goja"
	"QLToolsPro/utils/panel"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	time2 "time"

	"go.uber.org/zap"
)

// OnlineServer 获取在线服务
func OnlineServer() (res.ResCode, model.OnlineServer) {
	// 缓存中读取在线服务
	online, err := gcache.GetCache("online_server")
	if err != nil {
		zap.L().Debug("实时计算在线服务")
		// 实时计算
		var os model.OnlineServer
		// 获取启用服务器信息
		serverData := dao.GetPanelStartAllData()
		// 序列化数据
		data, err := json.Marshal(serverData)
		if err != nil {
			zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
			return res.CodeServerBusy, os
		}
		err = json.Unmarshal(data, &os.Online)
		if err != nil {
			zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
			return res.CodeServerBusy, os
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
					return res.CodeServerBusy, os
				}
				err = json.Unmarshal(eData, &os.Online[i].EnvData)
				if err != nil {
					zap.L().Error("[获取在线服务]失败，原因：" + err.Error())
					return res.CodeServerBusy, os
				}
			}

			// 获取面板已存在变量数量
			url := panel.StringHTTP(serverData[i].PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(serverData[i].PanelParams)
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
				return res.CodeServerBusy, os
			}

			// 判断返回状态
			if token.Code > 200 && token.Code < 500 {
				// 尝试获取授权
				go panel.GetPanelToken(serverData[i].PanelURL, serverData[i].PanelClientID, serverData[i].PanelClientSecret)

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
		}
		return res.CodeSuccess, os
	} else {
		// 返回缓存数据
		zap.L().Debug("缓存数据")
		zap.L().Debug(fmt.Sprintf("%s", online))
		var os model.OnlineServer
		err = json.Unmarshal([]byte(online.(string)), &os)
		if err != nil {
			zap.L().Error("[获取在线服务]错误，原因：" + err.Error())
			return res.CodeServerBusy, os
		}
		return res.CodeSuccess, os
	}
}

// OnlineUploadData 上传变量
func OnlineUploadData(uid any, p *model.OnlineEnvUpload) (res.ResCode, string) {
	var err error     // 错误变量
	var cache2 string // 临时缓存变量

	cache2 = p.EnvData

	// 不允许内容为空
	if p.EnvData == "" {
		return res.CodeEnvError, "禁止上传空内容"
	}

	var token model.PanelRes
	// 校验服务器ID
	serverData := dao.GetPanelIDData(p.ServerID)
	if serverData.ID == 0 {
		// 服务器不存在
		zap.L().Debug("[上传变量]提交容器不存在：" + strconv.Itoa(p.ServerID))
		return res.CodeEnvError, "提交服务器不在白名单"
	}

	// 校验变量名是否存在
	envData := dao.GetEnvNameData(p.EnvName)
	if envData.ID == 0 {
		// 变量不存在
		zap.L().Debug("[上传变量]提交变量名不存在：" + p.EnvName)
		return res.CodeEnvError, "提交变量不在白名单"
	}

	// 获取上传用户数据
	user := dao.GetUserIDData(uid)

	// 校验用户上传资格
	switch envData.EnvIsCharging {
	case 2:
		//	会员提交
		zap.L().Debug("会员提交")
		if time2.Now().Unix() > user.ActivationTime.Unix() {
			return res.CodeEnvError, "此项目需要会员资格提交"
		}
	case 3:
		//	积分提交
		zap.L().Debug("积分提交")
		if int(user.Integral)-envData.EnvNeedIntegral <= 0 {
			return res.CodeEnvError, "您的积分余额不足以提交此项目"
		}
	}

	// 转换切片
	envBind := strings.Split(serverData.EnvBinding, "@")
	// 校验变量是否处于容器白名单
	zap.L().Debug("[上传变量]校验变量是否处于容器白名单")
	zap.L().Debug("[上传变量]提交值：" + strconv.Itoa(int(envData.ID)))
	num := 0
	for i := 0; i < len(envBind); i++ {
		zap.L().Debug("[上传变量]变量绑定数据：" + envBind[i])
		if envBind[i] != "" {
			if envBind[i] == strconv.Itoa(int(envData.ID)) {
				num++
			}
		}
	}
	if num == 0 {
		return res.CodeEnvError, "变量不在容器白名单范围"
	}

	// 正则处理(检查是否符合规则)
	var s [][]string
	if envData.EnvRegex != "" {
		// 需要处理正则
		zap.L().Debug("需要处理正则")
		reg := regexp.MustCompile(envData.EnvRegex)
		// 匹配内容
		if reg != nil {
			s = reg.FindAllStringSubmatch(p.EnvData, -1)
			zap.L().Debug("[上传变量]" + fmt.Sprintf("%v", s))
			if len(s) == 0 {
				return res.CodeEnvError, "上传内容不符合规定"
			}
			cache2 = s[0][0]
		} else {
			return res.CodeServerBusy, ""
		}
	}

	// 校验变量配额
	code, edr, intNumber := CalculateQuantity(p.ServerID, envData.EnvMode, p.EnvName)
	if code == res.CodeServerBusy {
		zap.L().Debug("[上传变量]处理请求失败")
		return res.CodeServerBusy, ""
	} else if intNumber <= 0 {
		if envData.EnvMode != 3 {
			zap.L().Debug("[上传变量]限额已满，禁止提交")
			return res.CodeEnvError, "限额已满，禁止提交"
		}
	}

	// 检查重复提交
	var IsRepeat bool
	var EnvID int
	if IsRepeat, EnvID = CheckRepeat(edr, envData, cache2); IsRepeat != false {
		return res.CodeEnvError, "禁止提交重复数据"
	}

	// 是否启用插件
	if envData.EnvIsPlugin != false {
		var cache3 string
		var IsAdopt bool
		// 启用插件, 传入插件名称和变量
		zap.L().Debug("[上传变量]变量：" + envData.EnvName + "  启用插件：" + envData.EnvPluginName)
		IsAdopt, cache3 = goja.RunPlugin(envData.EnvPluginName, cache2)
		if IsAdopt != true {
			return res.CodeNoAdmittance, cache3
		}
		cache2 = cache3
	}

	// 提交到服务器
	var uploadData string // 上传数据
	var startDate string  // 上传后启用数据（更新模式）
	var IsNew bool        // 是否为新建模式
	url := panel.StringHTTP(serverData.PanelURL) + "/open/envs?t=" + strconv.Itoa(serverData.PanelParams)
	startDate = panel.StringHTTP(serverData.PanelURL) + "/open/envs/enable?t=" + strconv.Itoa(serverData.PanelParams)
	zap.L().Debug(url)

	// 指定上传数据
	if envData.EnvMode == 1 {
		// 新建模式
		zap.L().Debug("[上传变量]上传变量：新建模式")
		uploadData = `[{"value": "` + cache2 + `","name": "` + p.EnvName + `","remarks": "` + user.Username + `"}]`
	} else if envData.EnvMode == 2 {
		// 更新模式
		zap.L().Debug("[上传变量]上传变量：更新模式")
		/*
			1、获取传入变量的正则
			2、循环匹配正则
			3、匹配成功：更新、匹配失败：新建
		*/
		reg := regexp.MustCompile(envData.EnvUpdate)
		s3 := reg.FindAllStringSubmatch(cache2, -1)
		co := 0
		for i := 0; i < len(edr.Data); i++ {
			// 循环匹配正则, 判断面板变量名和传入变量名是否一致
			if edr.Data[i].Name == p.EnvName {
				// 一致, 获取变量正则部分
				reEnv := reg.FindAllStringSubmatch(edr.Data[i].Value, -1)
				// 判断匹配结果是否为空
				if len(reEnv) != 0 {
					// 判断两个正则值是否一致
					if reEnv[0][0] == s3[0][0] {
						// 一致，更新变量
						co = 0
						IsNew = false
						if serverData.PanelVersion {
							// 旧面板
							startDate = edr.Data[i].OldID
							uploadData = `{"_id": "` + edr.Data[i].OldID + `", "value": "` + cache2 + `","name": "` + p.EnvName + `","remarks": "` + user.Username + `"}`
						} else {
							// 新面板
							startDate = strconv.Itoa(edr.Data[i].ID)
							uploadData = `{"id": "` + strconv.Itoa(edr.Data[i].ID) + `", "value": "` + cache2 + `","name": "` + p.EnvName + `","remarks": "` + user.Username + `"}`
						}
						break
					} else {
						// 不一致，新建变量
						co++
					}
				}
			} else {
				// 面板没存在此变量
				co++
			}
		}

		if co != 0 {
			IsNew = true
			uploadData = `[{"value": "` + cache2 + `","name": "` + p.EnvName + `","remarks": "` + user.Username + `"}]`
		}
	} else {
		// 合并模式
		zap.L().Debug("上传变量：合并模式")
		if EnvID != -1 {
			vv := ""
			sList := strings.Split(edr.Data[EnvID].Value, "\n")
			if len(sList) != 1 {
				for _, str := range sList {
					vv += str + "\\n"
				}
				vv = vv + "\\n" + cache2
			} else {
				vv = edr.Data[EnvID].Value + envData.EnvMerge + cache2
			}
			if edr.Data[EnvID].OldID != "" {
				uploadData = `{"_id": "` + edr.Data[EnvID].OldID + `", "value": "` + vv + `","name": "` + p.EnvName + `","remarks": "` + envData.EnvRemarks + `"}`
			} else {
				uploadData = `{"id": "` + strconv.Itoa(edr.Data[EnvID].ID) + `", "value": "` + vv + `","name": "` + p.EnvName + `","remarks": "` + envData.EnvRemarks + `"}`
			}
		} else {
			uploadData = `[{"value": "` + cache2 + `","name": "` + p.EnvName + `","remarks": "` + envData.EnvRemarks + `"}]`
		}
	}

	var r []byte
	if envData.EnvMode == 1 {
		// 新建模式(POST)
		zap.L().Debug("新建模式(POST)")
		r, err = requests.Requests("POST", url, uploadData, serverData.PanelToken)
	} else if envData.EnvMode == 2 {
		// 更新模式(PUT)
		zap.L().Debug("更新模式(PUT)")
		if IsNew == false {
			r, err = requests.Requests("PUT", url, uploadData, serverData.PanelToken)
			// 启用禁用变量
			var EnableID string
			if serverData.PanelVersion {
				// 新版本
				EnableID = "[" + startDate + "]"
			} else {
				// 旧版本
				EnableID = `["` + startDate + `"]`
			}
			zap.L().Debug("启用禁用变量：" + EnableID)
			go func() {
				_, _ = requests.Requests("PUT", startDate, EnableID, serverData.PanelToken)
			}()
		} else {
			// 面板不存在变量时新建(POST)
			if intNumber <= 0 {
				zap.L().Debug("限额已满，禁止提交")
				return res.CodeEnvError, "限额已满，禁止提交"
			} else {
				r, err = requests.Requests("POST", url, uploadData, serverData.PanelToken)
			}
		}
	} else {
		// 合并模式(PUT)
		zap.L().Debug("【合并模式】上传内容：" + uploadData)
		if EnvID != -1 {
			r, err = requests.Requests("PUT", url, uploadData, serverData.PanelToken)
		} else {
			// 面板不存在合并模式变量时(POST)
			r, err = requests.Requests("POST", url, uploadData, serverData.PanelToken)
		}
	}

	if err != nil {
		return res.CodeServerBusy, ""
	}

	// 序列化内容
	if err = json.Unmarshal(r, &token); err != nil {
		zap.L().Error(err.Error())
		return res.CodeServerBusy, ""
	}

	if token.Code >= 400 && token.Code < 500 {
		// 尝试更新Token
		zap.L().Warn("上传错误警告：" + token.Message)
		go panel.GetPanelToken(serverData.PanelURL, serverData.PanelClientID, serverData.PanelClientSecret)
		return res.CodeEnvError, "发生一点小意外，请重新提交"
	} else if token.Code >= 500 {
		return res.CodeEnvError, "提交数据发生【500】错误，错误原因：" + token.Message
	}

	// 扣除积分
	if envData.EnvIsCharging == 3 {
		zap.L().Debug("扣除前：" + strconv.FormatInt(user.Integral, 10))
		user.Integral -= int64(envData.EnvNeedIntegral)
		zap.L().Debug("扣除后：" + strconv.FormatInt(user.Integral, 10))
		go dao.UpdateUserIntegral(user)
	}
	// 记录上传信息
	go dao.CreateRecordInfo(user.UserID, envData.EnvName, envData.EnvIsCharging, envData.EnvNeedIntegral)
	return res.CodeSuccess, ""
}

// CalculateQuantity 计算变量剩余位置
func CalculateQuantity(id, mode int, name string) (res.ResCode, model.EnvDataResult, int) {
	var token model.EnvDataResult
	// 获取变量数据
	count := dao.GetEnvNameCount(name)

	// 获取容器信息
	sData := dao.GetPanelDataByID(id)

	// 获取面板已存在变量数量
	url := panel.StringHTTP(sData.PanelURL) + "/open/envs?searchValue=&t=" + strconv.Itoa(sData.PanelParams)
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

// CheckRepeat 校验是否重复上传
func CheckRepeat(p model.EnvDataResult, data model.Env, env string) (bool, int) {
	var count = 0
	// 通过变量名获取上传模式
	if data.EnvMode == 1 {
		// 新建模式需要校验重复上传，更新模式无需
		count = 0
		for i := 0; i < len(p.Data); i++ {
			if p.Data[i].Value == env {
				count++
				break
			}
		}
		if count != 0 {
			return true, count
		}
	} else if data.EnvMode == 3 {
		// 合并模式，遍历所有表获取合并表
		count = 0
		if len(p.Data) == 0 {
			return false, count
		}
		for i := 0; i < len(p.Data); i++ {
			if p.Data[i].Name == data.EnvName {
				count = i
				break
			}
		}
		// 判断面板无此变量
		if count == 0 {
			return false, -1
		}

		// 根据分隔符处理面板上的数据
		var up = 0
		envList := strings.Split(p.Data[count].Value, data.EnvMerge)
		for i := 0; i < len(envList); i++ {
			if envList[i] == env {
				up++
				break
			}
		}
		if up != 0 {
			return true, count
		}
	}
	return false, count
}
