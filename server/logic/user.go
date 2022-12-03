// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 10:52
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : user.go

package logic

import (
	"QLToolsPro/server/dao"
	"QLToolsPro/server/gcache"
	"QLToolsPro/server/model"
	"QLToolsPro/utils/Sha1"
	"QLToolsPro/utils/jwt"
	"QLToolsPro/utils/requests"
	res "QLToolsPro/utils/response"
	"QLToolsPro/utils/snowflake"
	"QLToolsPro/utils/timeTools"
	"QLToolsPro/utils/wxpusher"
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image/color"
	"math/rand"
	"strconv"
	time2 "time"
	"unicode/utf8"
)

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// IP 地址序列化
type location struct {
	Country   string `json:"country"`
	ShortName string `json:"short_name"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Area      string `json:"area"`
	Isp       string `json:"isp"`
	Net       string `json:"net"`
	Ip        string `json:"ip"`
	Code      int    `json:"code"`
	Desc      string `json:"desc"`
}

// CaptMake 生成验证码
func CaptMake() (id, b64s string, err error) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, lerr := captcha.Generate()
	return lid, lb64s, lerr
}

// CaptVerify 验证captcha是否正确
func CaptVerify(id string, capt string) bool {
	if store.Verify(id, capt, false) {
		return true
	} else {
		return false
	}
}

// SignUp 注册业务
func SignUp(p *model.UserSignUp) (res.ResCode, string) {
	// 检查验证码是否正确
	if !CaptVerify(p.Id, p.Capt) {
		return res.CodeRegisterError, "验证码错误"
	}

	// 查询是否允许注册
	data, _ := dao.GetSetting("register")
	if data.Value != "1" {
		return res.CodeRegisterError, "用户注册已关闭"
	}

	// 判断是否已存在账户
	result, _ := dao.GetUserData(p.Email, p.Username)
	if result == true {
		return res.CodeRegisterError, "当前邮箱或用户名已被注册"
	}

	if utf8.RuneCountInString(p.Username) < 8 && utf8.RuneCountInString(p.Password) < 8 {
		return res.CodeRegisterError, "用户名或密码要求最低8位数"
	}

	// 密码加密
	p.Password = Sha1.Sha1(p.Password)

	// 生成用户UID
	userID := strconv.FormatInt(snowflake.GenID(), 10)

	// 构造User实例
	user := &model.User{
		UserID:   userID,
		Email:    p.Email,
		Username: p.Username,
		Password: p.Password,
		Integral: 0,
		IsAdmin:  false,
		IsState:  true,
	}

	// 判断是否为第一个账号
	c := dao.GetUserRecord()
	if c == 0 {
		// 第一个注册账号为管理员账号
		user.IsAdmin = true
	}

	// 添加赠送积分
	integral, _ := dao.GetSetting("regIntegral")
	i, _ := strconv.Atoi(integral.Value)
	user.Integral += int64(i)

	// 保存进数据库
	err := dao.InsertUser(user)
	if err != nil {
		zap.L().Error("Error inserting database, err:", zap.Error(err))
		return res.CodeServerBusy, "服务繁忙"
	}
	return res.CodeSuccess, "注册成功"
}

// SignIn 登录业务
func SignIn(p *model.UserSignIn, RemoteIP string) (res.ResCode, string) {
	// 检查验证码是否正确
	if !CaptVerify(p.Id, p.Capt) {
		return res.CodeLoginError, "验证码错误"
	}

	// 检查用户名是否存在
	result, user := dao.GetUserNameData(p.Username)
	if result == false {
		// 不存在
		return res.CodeLoginError, "账户不存在"
	} else {
		// 检查用户是否被封禁
		if user.IsState == false {
			return res.CodeLoginError, "账户已被管理员封禁"
		}

		// 邮箱存在,记录传入密码
		oPassword := p.Password

		// 判断密码是否正确
		if user.Password != Sha1.Sha1(oPassword) {
			return res.CodeLoginError, "密码错误"
		} else {
			// 密码正确, 校验是否异地登录
			if user.LoginIP != "" {
				b1, b2 := CheckIf(user.LoginIP, RemoteIP)
				if b2 {
					return res.CodeServerBusy, "服务繁忙"
				} else {
					if !b1 {
						// 异地登录
						return res.CodeAbnormalEnvironment, "账户登录环境异常"
					}
				}
			}

			// 密码正确, 返回生成的Token（userSecret：密码前六位）
			token, err := jwt.GenToken(user.UserID, user.Password[:6])
			if err != nil {
				zap.L().Error("An error occurred in token generation, err:", zap.Error(err))
				return res.CodeServerBusy, "服务繁忙"
			}

			// 记录登录IP
			go dao.UpdateUserLoginIP(RemoteIP, user.UserID)

			return res.CodeSuccess, token
		}
	}
}

// AppletLogin 小程序登录业务
func AppletLogin(p *model.AppletLogin) (res.ResCode, string) {
	// 获取OpenID
	appid := viper.GetString("wechat_applet.appid")
	secret := viper.GetString("wechat_applet.secret")

	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + p.Code + "&grant_type=authorization_code"
	data, err := requests.Requests("GET", url, "", "")
	if err != nil {
		zap.L().Error("[小程序登录]失败， 原因：" + err.Error())
		return res.CodeServerBusy, "登录失败"
	}

	//	序列化数据
	var app model.AppletRes
	err = json.Unmarshal(data, &app)
	if err != nil {
		zap.L().Error("[小程序登录]失败，原因：" + err.Error())
		return res.CodeServerBusy, "登录失败"
	}

	switch app.Errcode {
	case 40029:
		return res.CodeLoginError, "code 无效"
	case 45011:
		return res.CodeLoginError, "请求太频繁，请稍候再试"
	case 40226:
		return res.CodeLoginError, "高风险等级用户，小程序登录拦截"
	case -1:
		return res.CodeLoginError, "系统繁忙，请稍候再试"
	}

	// 状态正常， 数据库查询账号
	user := dao.GetUserAppOpenIDData(app.Openid)
	if user.AppOpenID == "" && user.UserID == "" {
		// 查询是否允许注册
		reg, _ := dao.GetSetting("register")
		if reg.Value != "1" {
			return res.CodeLoginError, "用户注册已关闭"
		}

		// 未注册账户，生成用户UID
		userID := strconv.FormatInt(snowflake.GenID(), 10)

		// 构建用户
		newUser := &model.User{
			UserID:    userID,
			AppOpenID: app.Openid,
			Username:  app.Openid,
			Password:  app.Openid,
			Integral:  0,
			IsAdmin:   false,
			IsState:   true,
		}

		// 添加赠送积分
		integral, _ := dao.GetSetting("regIntegral")
		i, _ := strconv.Atoi(integral.Value)
		newUser.Integral += int64(i)

		// 保存进数据库
		err = dao.InsertUser(newUser)
		if err != nil {
			zap.L().Error("Error inserting database, err:", zap.Error(err))
			return res.CodeLoginError, "创建账户失败"
		}

		// 再次登录
		user2 := dao.GetUserAppOpenIDData(app.Openid)
		token, err := jwt.GenToken(user2.UserID, user2.Password[:6])
		if err != nil {
			zap.L().Error("An error occurred in token generation, err:", zap.Error(err))
			return res.CodeServerBusy, "服务繁忙"
		}
		return res.CodeSuccess, token
	} else {
		// 注册账户，检查用户是否被封禁
		if user.IsState == false {
			return res.CodeLoginError, "账户已被管理员封禁"
		}
		token, err := jwt.GenToken(user.UserID, user.Password[:6])
		if err != nil {
			zap.L().Error("An error occurred in token generation, err:", zap.Error(err))
			return res.CodeServerBusy, "服务繁忙"
		}
		return res.CodeSuccess, token
	}
}

// GetUserOneData 用户信息：获取
func GetUserOneData(uid interface{}) (res.ResCode, model.UserData) {
	var UserData model.UserData
	data := dao.GetUserIDData(uid)
	UserData.UserID = data.UserID
	UserData.Email = data.Email
	UserData.Username = data.Username
	UserData.Integral = data.Integral
	UserData.UserWxpusher = data.UserWxpusher
	// 修改会员到期时间
	if time2.Now().Unix() > data.ActivationTime.Unix() {
		// 已到期
		UserData.IsActivationTime = false
		UserData.ActivationTime = "未开通或已过期"
	} else {
		// 未到期
		UserData.IsActivationTime = true
		UserData.ActivationTime = timeTools.SwitchTimeStampToDataYearNowTome(data.ActivationTime.Unix())
	}
	UserData.IsAdmin = data.IsAdmin
	return res.CodeSuccess, UserData
}

// UserSettings 获取一个配置信息
func UserSettings(name string) (model.WebSettings, res.ResCode) {
	data, err := dao.GetSetting(name)
	if err != nil {
		zap.L().Error(err.Error())
		return data, res.CodeServerBusy
	}

	// 限制前端只能获取公告信息
	if name == "register" {
		return data, res.CodeServerBusy
	}

	return data, res.CodeSuccess
}

// UserDivisionData 用户分页查询
func UserDivisionData(page, quantity string) (res.ResCode, model.UserPageData) {
	var data []model.User
	var pageData model.UserPageData
	var q int

	q, _ = strconv.Atoi(quantity)

	if page == "" {
		// 空值，默认获取前20条数据(第一页)
		data = dao.GetDivisionUserData(1, q)
	} else {
		// String转Int
		intPage, err := strconv.Atoi(page)
		if err != nil {
			// 类型转换失败，查询默认获取前20条数据(第一页)
			zap.L().Error(err.Error())
			data = dao.GetDivisionUserData(1, q)
		} else {
			// 查询指定页数的数据
			data = dao.GetDivisionUserData(intPage, q)
		}
	}

	// 查询总页数
	count := dao.GetUserRecord()
	// 计算页数
	z := count / int64(q)
	var y int64
	y = count % int64(q)

	if y != 0 {
		pageData.Page = z + 1
	} else {
		pageData.Page = z
	}

	// 删除密码
	for i := 0; i < len(data); i++ {
		if time2.Now().Unix() > data[i].ActivationTime.Unix() {
			data[i].ActivationTime = time2.Time{}
		}
		data[i].Password = "safety protection"
	}

	pageData.PageData = data

	return res.CodeSuccess, pageData
}

// UserSearch 用户名模糊查询
func UserSearch(s string) (res.ResCode, []model.User) {
	data := dao.GetUserNameSearch(s)
	// 删除密码
	for i := 0; i < len(data); i++ {
		data[i].Password = "safety protection"
	}
	return res.CodeSuccess, data
}

// UserInformationUpdate 用户数据更新
func UserInformationUpdate(p *model.UpdateUserData) res.ResCode {
	dao.UserInformationUpdate(p)
	return res.CodeSuccess
}

// UserInformationDelete 用户数据删除
func UserInformationDelete(p *model.DeleteUserData) res.ResCode {
	dao.UserInformationDelete(p)
	return res.CodeSuccess
}

// UserFindPwd 找回密码 - 发送Token
func UserFindPwd(p *model.UserFindPwd) (res.ResCode, string) {
	// 检查验证码是否正确
	if !CaptVerify(p.Id, p.Capt) {
		return res.CodeRePwdError, "验证码错误"
	}

	var result bool
	var user model.User
	// 判断是否已存在账户
	if p.FindType == "email" {
		result, user = dao.GetUserEmailData(p.Content)
		if result == false {
			return res.CodeRePwdError, "邮箱不存在"
		}
	} else {
		result, user = dao.GetUserNameData(p.Content)
		if result == false {
			return res.CodeRePwdError, "用户名不存在"
		}
	}

	// 判断是否绑定WxPusher
	if user.UserWxpusher == "" {
		return res.CodeRePwdError, "用户未绑定WxPusher，请联系管理员解决"
	}

	// 生成Token
	rand.Seed(time2.Now().Unix())
	num := rand.Intn(999999)
	value, _ := dao.GetSetting("web_title")
	data := "【" + value.Value + "】提示您。您的账户正在找回密码，请输入验证码：" + strconv.Itoa(num) + "。完成验证【验证码有效期：60分钟】"
	b := wxpusher.PluginSendMessageMass([]string{user.Username}, data)
	if b == false {
		return res.CodeRePwdError, "验证码发送失败，请稍后尝试"
	}

	// 发送成功，数据写入缓存
	var f model.CacheRecordPwd
	f.Code = strconv.Itoa(num)
	f.UserId = user.UserID
	v, err := json.Marshal(f)
	if err != nil {
		return res.CodeRePwdError, "验证码发送失败，请稍后尝试"
	}
	go gcache.TimingCache(strconv.Itoa(num), string(v), time2.Hour)

	// 返回
	return res.CodeSuccess, "验证码发送成功"
}

// UserRePwd 找回密码 - 修改密码
func UserRePwd(p *model.UserRePwd) (res.ResCode, string) {
	// 缓存查询Token
	code, err := gcache.GetCache(p.Code)
	if err != nil {
		return res.CodeRePwdError, "验证码已失效或不存在"
	}

	// 序列化内容
	var f model.CacheRecordPwd
	err = json.Unmarshal([]byte(code.(string)), &f)
	if err != nil {
		zap.L().Error("[找回密码 - 修改密码]错误，原因：" + err.Error())
		return res.CodeServerBusy, "服务繁忙"
	}

	// 修改密码
	err = dao.UpdateUserPwd(f.UserId, Sha1.Sha1(p.Password))
	if err != nil {
		zap.L().Error("[修改密码]错误, 原因：" + err.Error())
		return res.CodeRePwdError, "修改密码失败，请稍后重试"
	}

	// 删除Redis中的Token
	go gcache.DeleteCache(f.Code)

	return res.CodeSuccess, "修改密码成功"
}

// CheckIf 检查是否属于异地登录
func CheckIf(ip1, ip2 string) (bool, bool) {
	/*
		ip1：原始IP
		ip2：登录IP
		bool1：是否异地登录
		bool2：是否解析出错
	*/
	// IP地址相同，跳过验证
	if ip1 == ip2 {
		return true, false
	}

	// 查询IP地址
	url1 := "https://ip.useragentinfo.com/json?ip=" + ip1
	addr1, err := requests.Requests("GET", url1, "", "")
	if err != nil {
		return false, true
	}

	var l1 location
	var l2 location
	// 数据绑定
	err = json.Unmarshal(addr1, &l1)
	if err != nil {
		zap.L().Error(err.Error())
		return false, true
	}

	url2 := "https://ip.useragentinfo.com/json?ip=" + ip2
	addr2, err := requests.Requests("GET", url2, "", "")
	if err != nil {
		return false, true
	}
	// 数据绑定
	err = json.Unmarshal(addr2, &l2)
	if err != nil {
		zap.L().Error(err.Error())
		return false, true
	}

	if l1.Province == l2.Province {
		return true, false
	} else {
		return false, false
	}
}

// AbnormalEmail 登录异常 - 发送验证码
func AbnormalEmail(p *model.UserAbnormalEmail) (res.ResCode, string) {
	// 判断是否已存在账户
	result, user := dao.GetUserNameData(p.UserName)
	if result == false {
		return res.CodeAbnormalError, "用户不存在"
	}

	// 缓存查询Token
	_, err := gcache.GetCache(user.Username + "login")
	if err != nil {
		var f model.CacheRecordPwd

		// 生成验证码, 发送邮件
		rand.Seed(time2.Now().UnixNano())
		bytes := make([]byte, 5)
		for i := 0; i < 5; i++ {
			b := rand.Intn(26) + 65
			bytes[i] = byte(b)
		}
		zap.L().Debug("生成验证码：" + string(bytes))

		str := "您的登录验证码为：" + string(bytes) + "， (5分钟内有效，本邮件由系统自动发出，请勿直接回复)"
		zap.L().Debug("str地址：" + str)
		var uw []string
		uw = append(uw, user.UserWxpusher)
		b, msg := wxpusher.AdminSendMessage(uw, str)
		if !b {
			zap.L().Error("[WxPusher]发送失败，原因：" + msg)
			return res.CodeAbnormalError, "验证码发送失败，请稍等片刻再尝试"
		}

		// 发送成功，数据存入缓存
		f.Code = string(bytes)
		f.UserId = user.UserID
		v, err := json.Marshal(f)
		if err != nil {
			return res.CodeAbnormalError, "邮件发送失败，请稍等片刻再尝试"
		}
		go gcache.TimingCache(user.Username+"login", string(v), time2.Minute*5)

		// 返回
		return res.CodeSuccess, "验证码发送成功"
	}
	return res.CodeAbnormalError, "存在未过期验证码，请勿重复发送"
}

// AbnormalSignin 登录异常 - 登录
func AbnormalSignin(p *model.UserAbnormalSignin, RemoteIP string) (res.ResCode, string) {
	// 缓存查询Token
	uTk, err := gcache.GetCache(p.Username + "login")
	if err != nil {
		zap.L().Error("[登录异常-登录]失败，原因：" + err.Error())
		return res.CodeAbnormalError, "暂无登录验证码"
	}

	// 序列化字符串
	var f model.CacheRecordPwd
	err = json.Unmarshal([]byte(uTk.(string)), &f)
	if err != nil {
		zap.L().Error("[登录异常-登录]失败，原因：" + err.Error())
		return res.CodeServerBusy, "业务繁忙"
	}

	if p.VfCode != f.Code {
		return res.CodeAbnormalError, "验证码错误"
	}

	// 登录验证, 检查用户名是否存在
	result, user := dao.GetUserNameData(p.Username)
	if result == false {
		// 不存在
		return res.CodeLoginError, "用户名不存在"
	} else {
		// 检查用户是否被封禁
		if user.IsState == false {
			return res.CodeAbnormalError, "账户已被封禁"
		}

		// 邮箱存在,记录传入密码
		oPassword := p.Password

		// 判断密码是否正确
		if user.Password != Sha1.Sha1(oPassword) {
			return res.CodeAbnormalError, "密码错误"
		} else {
			// 密码正确, 返回生成的Token
			token, err := jwt.GenToken(user.UserID, user.Password[:6])
			if err != nil {
				zap.L().Error("An error occurred in token generation, err:", zap.Error(err))
				return res.CodeServerBusy, "服务繁忙"
			}

			// 记录登录IP
			go dao.UpdateUserLoginIP(RemoteIP, user.UserID)
			// 删除缓存中的验证码
			go gcache.DeleteCache(p.Username + "login")

			return res.CodeSuccess, token
		}
	}
}
