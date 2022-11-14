package logic

import (
	_const "QLToolsPro/server/const"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
)

// RemoteVersion 系统远程信息
type RemoteVersion struct {
	LocVersion string
	Version    string
	Notice     string
}

// CheckVersion 检查版本更新
func CheckVersion() (RemoteVersion, res.ResCode) {
	// 版本号
	var v RemoteVersion

	// 获取仓库版本信息
	url := "https://version.6b7.xyz/qltoolspro_version.json"
	r, err := requests.Requests("GET", url, "", "")
	if err != nil {
		return v, res.CodeServerBusy
	}
	// 序列化内容
	err = json.Unmarshal(r, &v)
	if err != nil {
		zap.L().Error(err.Error())
		return v, res.CodeServerBusy
	}
	v.LocVersion = _const.LocVersion

	return v, res.CodeSuccess
}
