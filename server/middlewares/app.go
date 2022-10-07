package middlewares

import (
	_const "QLToolsPro/server/const"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	time2 "time"
)

var AppLicence bool

// AppAuth 小程序授权认证中间件
func AppAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if AppLicence {
			c.Next()
		} else {
			res.ResErrorWithMsg(c, res.CodeServerBusy, "未获取授权")
			c.Abort()
			return
		}
	}
}

// LicenseCheck 检查运行许可
func LicenseCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		if time2.Now().Unix() < _const.LicenseDate {
			c.Next()
		} else {
			// 许可到期
			zap.L().Warn("许可失效，停止服务")
			res.ResError(c, res.CodeServerBusy)
			c.Abort()
			return
		}
	}
}

// LoginRes 许可证返回数据
type LoginRes struct {
	Code int  `json:"code"`
	Data bool `json:"data"`
}

// CheckAppLicence 小程序授权检查
func CheckAppLicence() {
	AppLicence = false
	var LR LoginRes
	// 检查授权
	url := "https://licence.pro.6b7.org/v1/api/user/applet/get"
	data := `{"user_name": ` + viper.GetString("licence.username") + `}`
	ResData, err := requests.Requests("POST", url, data, "")
	if err != nil {
		AppLicence = false
	}
	if err = json.Unmarshal(ResData, &LR); err != nil {
		AppLicence = false
	}
	if LR.Data {
		AppLicence = true
	} else {
		AppLicence = false
	}
}
