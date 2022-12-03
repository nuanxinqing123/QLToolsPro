// -*- coding: utf-8 -*-
// @Time    : 2022/8/26 10:12
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : cdk.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"bufio"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

// CDKEYDivisionCDKData CD-KEY分页查询
func CDKEYDivisionCDKData(ctype, page, quantity string) (res.ResCode, model.CDKPageData) {
	var data []model.CDK
	var pageData model.CDKPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.CDKEYDivisionCDKData(ctype, 1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.CDKEYDivisionCDKData(ctype, 1, q)
		} else {
			// 查询指定页数的数据
			data = dao.CDKEYDivisionCDKData(ctype, intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetCDKEYDataPage(ctype)
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

// CDKEYSearch CD-KEY数据查询
func CDKEYSearch(ctype, s string) (res.ResCode, model.CDK) {
	c := dao.CDKEYSearch(ctype, s)
	if c.CdKey == "" {
		return res.CodeCDKError, c
	}
	return res.CodeSuccess, c
}

// CDKEYAdd 批量生成CD-KEY
func CDKEYAdd(p *model.CreateCDK) res.ResCode {
	// 判断本地是否还有遗留文件
	_, err := os.Stat("CDK.txt")
	if err == nil {
		// 删除旧文件
		err = os.Remove("CDK.txt")
		if err != nil {
			zap.L().Error(err.Error())
			return res.CodeServerBusy
		}
	}

	// 创建记录数组
	var li []string

	// 创建对象
	cdk := new(model.CDK)

	if p.CdKeyType == "integral" {
		cdk.CdKeyType = "integral"
		cdk.CdKeyIntegral = p.CdKeyIntegral
	} else {
		cdk.CdKeyType = "vip"
		cdk.CdKeyValidityPeriod = int64(p.CdKeyValidityPeriod)
	}

	cdk.CdKeyState = true
	cdk.CdKeyRemarks = p.CdKeyRemarks
	// 获取生成数量
	for i := 0; i < p.CdKeyCount; i++ {
		// 生成用户UID
		uid := ksuid.New()
		cdk.CdKey = p.CdKeyPrefix + uid.String()

		// 加入数组
		li = append(li, cdk.CdKey)

		// 写入数据库
		dao.CDKEYAdd(cdk)
	}

	// 创建CDK.txt并写入数据
	filepath := "CDK.txt"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		zap.L().Error(err.Error())
		return res.CodeCDKError
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			zap.L().Error(err.Error())
		}
	}(file)

	// 写入CDK数据
	writer := bufio.NewWriter(file)
	for i := 0; i < len(li); i++ {
		_, err2 := writer.WriteString(li[i] + "\n")
		if err2 != nil {
			zap.L().Error(err2.Error())
			return res.CodeCDKError
		}
	}
	err = writer.Flush()
	if err != nil {
		zap.L().Error(err.Error())
		return res.CodeCDKError
	}

	return res.CodeSuccess
}

// CDKEYRemarksSearch CD-KEY标识查询
func CDKEYRemarksSearch(s string) (res.ResCode, []model.CDK) {
	return res.CodeSuccess, dao.CDKEYRemarksSearch(s)
}

// CDKEYDataDelete CD-KEY删除本地数据
func CDKEYDataDelete() {
	time.Sleep(time.Second * 3)
	err := os.Remove("CDK.txt")
	if err != nil {
		zap.L().Error(err.Error())
	}
}

// CDKEYBatchUpdate 修改CD-KEY
func CDKEYBatchUpdate(p *model.UpdateCDK) res.ResCode {
	// 更新CD-KEY数据
	if err := dao.CDKEYBatchUpdate(p); err != nil {
		zap.L().Error("Error update database, err:", zap.Error(err))
		return res.CodeEnvError
	}
	return res.CodeSuccess
}

// CDKEYDelete 删除CD-KEY
func CDKEYDelete(p *model.DelCDK) res.ResCode {
	// 删除CD-KEY数据
	if err := dao.CDKEYDelete(p); err != nil {
		zap.L().Error("Error update database, err:", zap.Error(err))
		return res.CodeEnvError
	}
	return res.CodeSuccess
}
