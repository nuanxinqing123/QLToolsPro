package license

import (
	_const "QLToolsPro/server/const"
	"QLToolsPro/utils/Sha1"
	"QLToolsPro/utils/requests"
	jsoniter "github.com/json-iterator/go"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
	"github.com/super-l/machine-code/machine"
	"go.uber.org/zap"
	"regexp"
	time2 "time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

/* 网络验证模块 */

const AppKey = "59faaba418bd60d4e61b66c636328789"

// LoginRes 许可证返回数据
type LoginRes struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Id        string `json:"id"`        // 卡密在整个平台的排序ID
	EndTime   string `json:"end_time"`  // 卡密到期时间
	Amount    string `json:"amount"`    // 卡密生成时选择的时长
	Available int64  `json:"available"` // 卡密还剩多少秒到期
	Token     string `json:"token"`     // 签名
	Date      string `json:"date"`      // 服务器时间(YmdHM)
	Time      int64  `json:"time"`      // 服务器时间戳
	Imei      string `json:"imei"`      // 卡密绑定的设备码
	StateCode string `json:"statecode"` // 心跳状态码
}

// LoginLicense 许可证登录
func LoginLicense() (bool, string) {
	var LR LoginRes
	username := viper.GetString("licence.username")
	password := viper.GetString("licence.password")
	aqCode := viper.GetString("licence.aqCode")
	url := "https://w.t3yanzheng.com/39A83D9147AD81B5?user=" + username + "&pass=" + password + "&imei=" + GetMachineCode() + "&s=" + EncryptionSign()
	ResData, err := requests.Requests("POST", url, "", "")
	if err != nil {
		url2 := "https://w.t3data.net/39A83D9147AD81B5?user=" + username + "&pass=" + password + "&imei=" + GetMachineCode() + "&s=" + EncryptionSign()
		ResData, _ = requests.Requests("POST", url2, "", "")
		if err = json.Unmarshal(ResData, &LR); err != nil {
			zap.L().Error("[登录]失败，原因：" + err.Error())
			return false, "许可证验证失败"
		}
	} else {
		if err = json.Unmarshal(ResData, &LR); err != nil {
			zap.L().Error("[登录]失败，原因：" + err.Error())
			return false, "许可证验证失败"
		}
	}

	if LR.Code != 200 {
		return false, LR.Msg
	}
	_const.LicenseDate = time2.Now().Unix() + LR.Available

	// 用户信息发送至服务器【用户名、安全码、到期时间】
	go UpdateServerData(username, aqCode, LR.EndTime)

	return true, LR.EndTime
}

// GetMachineCode 获取机器码
func GetMachineCode() string {
	// 系统版本号 + PlatformUUID
	// 签名：Sha1 + RC4
	version, _ := host.KernelVersion()
	uuid, _ := machine.GetPlatformUUID()
	return Sha1.DESPlus(Sha1.Sha1(version + uuid))[:32]
}

// EncryptionSign 请求签名
func EncryptionSign() string {
	// 取MD5值("kami=" + 许可证 + "&imei=" + 机器码 + "&t=" + 取现行时间戳(2) + "&" + APPKEY)
	username := viper.GetString("licence.username")
	password := viper.GetString("licence.password")
	MachineCode := GetMachineCode()
	return Sha1.Md5("user=" + username + "&pass=" + password + "&imei=" + MachineCode + "&" + AppKey)
}

// GetIPaddress 获取IP地址
func GetIPaddress() string {
	// 查询IP地址
	url := "http://23.80.5.90/ip.php"
	addr, _ := requests.Requests("GET", url, "", "")

	// 正则匹配IP地址
	reg := regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	s := reg.FindAllStringSubmatch(string(addr), -1)
	zap.L().Debug("IP地址：" + s[0][0])
	return s[0][0]
}

// UpdateServerData 修改服务器储存数据
func UpdateServerData(u, a, e string) {
	url := "https://licence.pro.6b7.org/v1/api/user/data/update/admin/update"
	data := `{"user_name": "` + u + `", "aq_code": "` + a + `", "due_date": "` + e + `"}`
	_, _ = requests.Requests("PUT", url, data, "")
}
