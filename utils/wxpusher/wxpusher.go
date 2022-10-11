// -*- coding: utf-8 -*-
// @Time    : 2022/8/19 8:45
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : wxpusher.go

package wxpusher

import (
	"QLToolsPro/server/dao"
	WxModel "QLToolsPro/utils/wxpusher/model"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

/* 封装WxPusher方法 */

// GetWxPusherQRCode 获取订阅二维码
func GetWxPusherQRCode() (WxModel.CreateQrcodeResult, error) {
	title, _ := dao.GetSetting("web_title")
	qrcode := WxModel.Qrcode{AppToken: viper.GetString("wxpusher.wp_app_token"), Extra: title.Value + "订阅用户"}
	return CreateQrcode(&qrcode)
}

// PluginSendMessageMass 插件发送消息群发
// UserName：用户名 数组
// Msg：消息内容
func PluginSendMessageMass(UserName []string, Msg string) bool {
	var UserWxpusher []string

	for i := 0; i < len(UserName); i++ {
		// UserName查询用户WxPusherID
		_, user := dao.GetUserNameData(UserName[i])
		if user.UserWxpusher != "" {
			UserWxpusher = append(UserWxpusher, user.UserWxpusher)
		}
	}

	msg := WxModel.NewMessage(viper.GetString("wxpusher.wp_app_token")).SetContent(Msg).AddUIdMass(UserWxpusher)
	msgArr, err := SendMessage(msg)
	zap.L().Debug("MsgArr Return:" + fmt.Sprintf("%v", msgArr))
	if err != nil {
		// 发送失败
		zap.L().Error("WxPusher Send Message Error:" + err.Error())
		return false
	} else {
		// 发送成功
		zap.Any("WxPusher Return:", msgArr)
		return true
	}
}

// AdminSendMessage 管理员发送消息群发
// UserWxpusher：用户WxpusherID 数组
// Msg：消息内容
func AdminSendMessage(UserWxpusher []string, Msg string) (bool, string) {
	msg := WxModel.NewMessage(viper.GetString("wxpusher.wp_app_token")).SetContent(Msg).AddUIdMass(UserWxpusher)
	msgArr, err := SendMessage(msg)
	if err != nil {
		// 发送失败
		zap.L().Error("WxPusher Send Message Error:" + err.Error())
		return false, "发送失败"
	} else {
		// 发送成功
		zap.Any("WxPusher Return:", msgArr)
		return true, "发送失败"
	}
}
