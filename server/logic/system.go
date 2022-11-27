package logic

import (
	_const "QLToolsPro/server/const"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"go.uber.org/zap"
)

// CheckVersion 检查版本更新
func CheckVersion() (model.RemoteVersion, res.ResCode) {
	// 版本号
	var v model.RemoteVersion

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

// SystemSoftwareUpdate 更新软件
//func SystemSoftwareUpdate(p *model.SoftWareGOOS) (res.ResCode, string) {
//	if runtime.GOOS == "windows" {
//		return res.CodeSystemError, "Windows系统不支持此功能"
//	}
//	// 获取版本号
//	var v model.Ver
//	url := "https://version.6b7.xyz/qltools_version.json"
//	r, _ := requests.Requests("GET", url, "", "")
//	_ = json.Unmarshal(r, &v)
//	if v.Version == _const.Version {
//		return res.CodeSystemError, "已经是最新版本"
//	}
//
//	// 更新程序
//	go UpdateSoftWare(v.Version, p.Framework)
//
//	return res.CodeSuccess, "程序已进入自动更新任务，如果更新失败请手动更新"
//}

//func UpdateSoftWare(version, GOOS string) {
//	// 更新地址
//	var url string
//	url = "https://github.com/nuanxinqing123/QLTools/releases/download/" + version
//
//	if GOOS == "amd64" {
//		url += "/QLTools-linux-amd64"
//	} else if GOOS == "arm64" {
//		url += "/QLTools-linux-arm64"
//	} else {
//		url += "/QLTools-linux-arm"
//	}
//	zap.L().Debug("Download: " + url)
//
//	err := doUpdate(url)
//	if err != nil {
//		zap.L().Error(err.Error())
//	}
//}

//func doUpdate(url string) error {
//	resp, err := requests.Down(url)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	err = update.Apply(resp.Body, update.Options{})
//	if err != nil {
//		// error handling
//		return err
//	}
//	return nil
//}
