// -*- coding: utf-8 -*-
// @Time    : 2022/8/31 10:25
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : recharge.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"time"
	"unicode"
)

// UserRechargeIntegral 用户充值：用户积分充值
func UserRechargeIntegral(uid any, p *model.UserRecharge) (res.ResCode, string) {
	// 校验CDK是否存在
	cdkData := dao.CDKEYDataSearch(p.RechargeCDK)
	if cdkData.CdKey == "" {
		return res.CodeCDKError, "请检查您的CDK是否有效"
	} else if cdkData.CdKeyState == false {
		return res.CodeCDKError, "CDK已被使用或已失效"
	} else if cdkData.CdKeyType != "integral" {
		return res.CodeCDKError, "充值卡密类型错误"
	}

	// 获取用户数据
	userData := dao.GetUserIDData(uid)
	// 充值用户额度
	userData.Integral += int64(cdkData.CdKeyIntegral)
	dao.UpdateUserIntegral(userData)

	// 已使用,禁用CDK
	cdkData.CdKeyState = false
	go dao.UpdateFalseCDK(cdkData)
	// 记录充值记录
	rechargeLog := new(model.Recharge)
	rechargeLog.RechargeUID = uid.(string)
	rechargeLog.RechargeType = "积分"
	rechargeLog.RechargeCDK = cdkData.CdKey
	go dao.InsertUserRechargeLog(rechargeLog)
	return res.CodeSuccess, "充值成功"
}

// UserRechargeVIP 用户充值：用户会员充值
func UserRechargeVIP(uid any, p *model.UserRecharge) (res.ResCode, string) {
	// 校验CDK是否存在
	cdkData := dao.CDKEYDataSearch(p.RechargeCDK)
	if cdkData.CdKey == "" {
		return res.CodeCDKError, "请检查您的CDK是否有效"
	} else if cdkData.CdKeyState == false {
		return res.CodeCDKError, "CDK已被使用或已失效"
	} else if cdkData.CdKeyType != "vip" {
		return res.CodeCDKError, "充值卡密类型错误"
	}

	// 获取用户数据
	userData := dao.GetUserIDData(uid)
	// 充值用户额度
	zap.L().Debug("初始额度：" + fmt.Sprintf("%s", userData.ActivationTime))
	if time.Now().Unix() > userData.ActivationTime.Unix() {
		switch cdkData.CdKeyValidityPeriod {
		case 1:
			userData.ActivationTime = time.Now().AddDate(0, 0, 1)
		case 7:
			userData.ActivationTime = time.Now().AddDate(0, 0, 7)
		case 14:
			userData.ActivationTime = time.Now().AddDate(0, 0, 14)
		case 31:
			userData.ActivationTime = time.Now().AddDate(0, 1, 0)
		case 92:
			userData.ActivationTime = time.Now().AddDate(0, 3, 0)
		case 183:
			userData.ActivationTime = time.Now().AddDate(0, 6, 0)
		case 365:
			userData.ActivationTime = time.Now().AddDate(1, 0, 0)
		}
	} else {
		switch cdkData.CdKeyValidityPeriod {
		case 1:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 1)
		case 7:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 7)
		case 14:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 14)
		case 31:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 1, 0)
		case 92:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 3, 0)
		case 183:
			userData.ActivationTime = userData.ActivationTime.AddDate(0, 6, 0)
		case 365:
			userData.ActivationTime = userData.ActivationTime.AddDate(1, 0, 0)
		}
	}

	// 更新用户数据
	zap.L().Debug("充值后额度：" + fmt.Sprintf("%s", userData.ActivationTime))
	dao.UpdateUserDataSave(userData)

	// 已使用,禁用CDK
	cdkData.CdKeyState = false
	go dao.UpdateFalseCDK(cdkData)
	// 记录充值记录
	rechargeLog := new(model.Recharge)
	rechargeLog.RechargeUID = uid.(string)
	rechargeLog.RechargeType = "会员"
	rechargeLog.RechargeCDK = cdkData.CdKey
	go dao.InsertUserRechargeLog(rechargeLog)
	return res.CodeSuccess, "充值成功"
}

// RechargeUserDivisionData 用户充值：记录查询
func RechargeUserDivisionData(uid any, page, quantity string) (res.ResCode, model.RechargePage) {
	var recharge []model.Recharge
	var rechargePage model.RechargePage
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		recharge = dao.RechargeUserDivisionData(uid, 1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			recharge = dao.RechargeUserDivisionData(uid, 1, q)
		} else {
			// 查询指定页数的数据
			recharge = dao.RechargeUserDivisionData(uid, intPage, q)
		}
	}

	// 查询总页数
	count := dao.RechargeUserCount(uid)
	// 计算页数
	z := count / int64(q)
	var y int64
	y = count % int64(q)

	if y != 0 {
		rechargePage.Page = z + 1
	} else {
		rechargePage.Page = z
	}
	rechargePage.PageData = recharge

	return res.CodeSuccess, rechargePage
}

// RechargeDivisionData 充值记录分页查询
func RechargeDivisionData(page, quantity string) (res.ResCode, model.RechargePage) {
	var data []model.Recharge
	var pageData model.RechargePage
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.RechargeDivisionData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.RechargeDivisionData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.RechargeDivisionData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetRechargeDataPage()
	// 计算页数
	z := count / int64(q)
	var y int64
	y = count % int64(q)

	if y != 0 {
		pageData.Page = z + 1
	} else {
		pageData.Page = z
	}
	pageData.PageData = data

	return res.CodeSuccess, pageData
}

// RechargeSearch Recharge数据查询
func RechargeSearch(s string) (res.ResCode, []model.Recharge) {
	// 判断是CDK搜索还是UserID搜索
	IsCDK := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			IsCDK = true
			break
		}
	}

	return res.CodeSuccess, dao.RechargeSearch(s, IsCDK)
}

// CDKEYUserRechargeIntegral 指定用户充值(积分)
func CDKEYUserRechargeIntegral(p *model.AdminRecharge) (res.ResCode, string) {
	// 获取用户数据
	userData := dao.GetUserIDData(p.UserID)
	if userData.Username == "" && userData.UserID == "" {
		return res.CodeCDKError, "充值账户不存在"
	}

	if p.RechargeType == 1 {
		// 充值用户额度
		userData.Integral += int64(p.RechargeNumber)
		dao.UpdateUserIntegral(userData)

		// 记录充值记录
		rechargeLog := new(model.Recharge)
		rechargeLog.RechargeUID = userData.UserID
		rechargeLog.RechargeType = "积分"
		rechargeLog.RechargeCDK = "管理员充值"
		go dao.InsertUserRechargeLog(rechargeLog)
		return res.CodeSuccess, "充值成功"
	} else {
		// 充值用户额度
		zap.L().Debug("初始额度：" + fmt.Sprintf("%s", userData.ActivationTime))
		if time.Now().Unix() > userData.ActivationTime.Unix() {
			switch p.RechargeNumber {
			case 1:
				userData.ActivationTime = time.Now().AddDate(0, 0, 1)
			case 7:
				userData.ActivationTime = time.Now().AddDate(0, 0, 7)
			case 14:
				userData.ActivationTime = time.Now().AddDate(0, 0, 14)
			case 31:
				userData.ActivationTime = time.Now().AddDate(0, 1, 0)
			case 92:
				userData.ActivationTime = time.Now().AddDate(0, 3, 0)
			case 183:
				userData.ActivationTime = time.Now().AddDate(0, 6, 0)
			case 365:
				userData.ActivationTime = time.Now().AddDate(1, 0, 0)
			}
		} else {
			switch p.RechargeNumber {
			case 1:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 1)
			case 7:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 7)
			case 14:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 0, 14)
			case 31:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 1, 0)
			case 92:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 3, 0)
			case 183:
				userData.ActivationTime = userData.ActivationTime.AddDate(0, 6, 0)
			case 365:
				userData.ActivationTime = userData.ActivationTime.AddDate(1, 0, 0)
			}
		}
		// 更新用户数据
		zap.L().Debug("充值后额度：" + fmt.Sprintf("%s", userData.ActivationTime))
		dao.UpdateUserDataSave(userData)

		// 记录充值记录
		rechargeLog := new(model.Recharge)
		rechargeLog.RechargeUID = userData.UserID
		rechargeLog.RechargeType = "会员"
		rechargeLog.RechargeCDK = "管理员充值"
		go dao.InsertUserRechargeLog(rechargeLog)
		return res.CodeSuccess, "充值成功"
	}
}
