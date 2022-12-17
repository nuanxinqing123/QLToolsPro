package logic

import (
	_const "QLToolsPro/server/const"
	"QLToolsPro/server/dao"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"github.com/staktrace/go-update"
	"go.uber.org/zap"
	"runtime"
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
func SystemSoftwareUpdate(p *model.SoftWareGOOS) (res.ResCode, string) {
	if runtime.GOOS == "windows" {
		return res.CodeSystemError, "Windows系统不支持此功能"
	}
	// 获取版本号
	var v model.RemoteVersion
	url := "https://version.6b7.xyz/qltoolspro_version.json"
	r, err := requests.Requests("GET", url, "", "")
	if err != nil {
		zap.L().Error("[系统更新]：" + err.Error())
		return res.CodeSystemError, "系统更新发生错误，详细请查看日志"
	}
	if err = json.Unmarshal(r, &v); err != nil {
		zap.L().Error("[系统更新]：" + err.Error())
		return res.CodeSystemError, "系统更新发生错误，详细请查看日志"
	}
	if v.Version == _const.LocVersion {
		return res.CodeSystemError, "已经是最新版本"
	}

	// 更新程序
	go UpdateSoftWare(v.Version, p.Framework)

	return res.CodeSuccess, "已开始自动更新（完成后需要手动重启），如果更新失败请手动更新"
}

// UpdateSoftWare 更新程序
func UpdateSoftWare(version, GOOS string) {
	// 更新地址
	var url string
	url = "https://version.6b7.xyz/update/QLToolsPro/" + version

	if GOOS == "amd64" {
		url += "/QLToolsPro-linux-amd64"
	} else {
		url += "/QLToolsPro-linux-arm64"
	}
	zap.L().Debug("Download: " + url)

	err := doUpdate(url)
	if err != nil {
		zap.L().Error("[系统更新]：" + err.Error())
		return
	}

	// Kill Main
	//go func() {
	//	// 等待两秒钟
	//	time.Sleep(time.Second * 2)
	//
	//	// Kill
	//	zap.L().Debug("进程PID：" + strconv.Itoa(syscall.Getpid()))
	//	cmd := exec.Command("/bin/bash", "-c", "kill -SIGHUP "+strconv.Itoa(syscall.Getpid()))
	//
	//	// 调用命令，如果发生错误，记录错误信息
	//	err = cmd.Start()
	//	if err != nil {
	//		zap.L().Error("[重启]：" + err.Error())
	//	}
	//}()
}

// 更新覆盖源文件
func doUpdate(url string) error {
	resp, err := requests.Down(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
		return err
	}
	return nil
}

// TaskDataQuery 查询任务
func TaskDataQuery(tp string) model.Task {
	return dao.GetCronTask(tp)
}
