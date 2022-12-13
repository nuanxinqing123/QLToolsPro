package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
	"strconv"
)

// CronTaskDivisionData 分页查询
func CronTaskDivisionData(page, quantity string) (res.ResCode, model.TaskPageData) {
	var data []model.Task
	var pageData model.TaskPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.CronTaskDivisionData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.CronTaskDivisionData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.CronTaskDivisionData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.CronTaskDivisionDataPage()
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
