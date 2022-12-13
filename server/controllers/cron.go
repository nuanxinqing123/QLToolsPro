package controllers

import (
	"QLToolsPro/server/logic"
	res "QLToolsPro/utils/response"
	"github.com/gin-gonic/gin"
)

// CronTaskDivisionData 分页查询
func CronTaskDivisionData(c *gin.Context) {
	// 获取查询页码
	page := c.Query("page")
	quantity := c.Query("quantity")
	resCode, data := logic.CronTaskDivisionData(page, quantity)

	switch resCode {
	case res.CodeSuccess:
		// 查询成功
		res.ResSuccess(c, data)
	}
}

// CronTaskAdd 创建任务
//func CronTaskAdd(c *gin.Context) {
//
//}
