// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:26
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : time.go

package timeTools

import "time"

// SwitchTimeStampToDataYear 将传入的时间戳转为时间
func SwitchTimeStampToDataYear(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02")
}

// SwitchTimeStampToDataYearNowTome 将传入的时间戳转为时间
func SwitchTimeStampToDataYearNowTome(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
