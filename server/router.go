// -*- coding: utf-8 -*-
// @Time    : 2022/8/17 21:07
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : router.go

package server

import (
	"QLToolsPro/server/controllers"
	"QLToolsPro/server/logger"
	"QLToolsPro/server/middlewares"
	"QLToolsPro/static/bindata"
	"html/template"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Setup() *gin.Engine {
	// 创建服务
	r := gin.New()

	// 配置中间件
	{
		// 配置日志
		if viper.GetString("app.mode") == "" {
			r.Use(logger.GinLogger(), logger.GinRecovery(true))
		}
	}

	// 前端静态文件
	{
		// 加载模板文件
		t, err := loadTemplate()
		if err != nil {
			panic(err)
		}
		r.SetHTMLTemplate(t)

		// 加载静态文件
		fs := assetfs.AssetFS{
			Asset:     bindata.Asset,
			AssetDir:  bindata.AssetDir,
			AssetInfo: nil,
			Prefix:    "assets",
		}
		r.StaticFS("/static", &fs)

		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	}

	// 路由组
	{
		// 开放权限组
		open := r.Group("/v1/api")
		// 配置跨域
		open.Use(cors.Default())
		{
			// 生成验证码
			open.GET("verification/code", controllers.CreateVerificationCode)
			// 账户注册
			open.POST("user/signup", controllers.SignUpHandle)
			// 账户登录
			open.POST("user/signin", controllers.SignInHandle)
			// 小程序登录
			open.POST("user/applet/login", controllers.AppletLoginHandle)
			// 登录异常-发送验证码
			open.POST("user/abnormal/code", controllers.UserAbnormalCode)
			// 登录异常-登录
			open.POST("user/abnormal/signin", controllers.UserAbnormalSignin)
			// 找回密码 - 发送验证码
			open.POST("findpwd/message", controllers.UserFindPwd)
			// 找回密码 - 修改密码˚Ω
			open.PUT("findpwd/repwd", controllers.UserRePwd)
			// 检查Token是否有效
			open.POST("check/token", middlewares.UserAuth(), controllers.CheckToken)
			// 检查管理员 Token是否有效
			open.POST("check/token/admin", middlewares.AdminAuth(), controllers.CheckToken)
			// 网站设置获取
			open.GET("set/settings", controllers.GetSettings)
			// 系统版本信息
			open.GET("system/version", controllers.SystemVersion)
		}

		// 用户权限组
		user := r.Group("/v2/api")
		user.Use(middlewares.UserAuth())
		// 配置跨域
		user.Use(cors.Default())
		{
			// 用户首页
			{
				// 用户信息获取
				user.GET("user/data", controllers.GetUserOneData)
				// 检查用户WxPusher订阅状态
				user.GET("user/wxpusher/state", controllers.UserWxpusherState)
				// 获取WxPusher订阅二维码
				user.GET("user/wxpusher/qrcode", controllers.UserWxpusherQrcode)
				// 更新用户WxPusherID
				user.PUT("user/wxpusher/update", controllers.UserWxpusherUpdate)
				// 公告获取
				user.GET("user/settings", controllers.UserSettings)
			}

			// 变量提交
			{
				// 在线服务
				user.GET("online/server", controllers.OnlineServer)
				// 上传内容
				user.POST("online/upload/data", controllers.OnlineUploadData)
			}

			// 变量管理
			//{
			//	// 分页查询
			//	// 筛选查询
			//	// 用户修改变量
			//	// 用户删除变量
			//}

			// 充值服务
			{
				// 积分充值
				user.POST("user/recharge/integral", controllers.UserRechargeIntegral)
				// 会员充值
				user.POST("user/recharge/vip", controllers.UserRechargeVIP)
			}

			// 上传记录
			{
				// 分页查询
				user.GET("record/user/division/data", controllers.RecordUserDivisionData)
			}

			// 充值记录
			{
				// 分页查询
				user.GET("consumption/user/division/data", controllers.RechargeUserDivisionData)
			}
		}

		// 管理员权限组
		admin := r.Group("/v3/api")
		admin.Use(middlewares.AdminAuth())
		// 配置跨域
		admin.Use(cors.Default())
		{
			// 管理员首页
			{
				// 首页数据
				admin.GET("index/data", controllers.AdminIndexData)
				// 更新系统
				admin.POST("system/software/update", controllers.SystemSoftwareUpdate)
				// 关闭/重启系统
				admin.GET("system/state", controllers.SystemState)
				// 查询任务
				admin.GET("task/data/query", controllers.TaskDataQuery)
			}

			// 面板管理
			{
				// 面板分页查询
				admin.GET("panel/division/data", controllers.PanelDivisionData)
				// 获取面板简易数据
				admin.GET("panel/division/data/simple", controllers.PanelDivisionDataSimple)
				// 面板新增
				admin.POST("panel/add", controllers.PanelAdd)
				// 面板修改
				admin.PUT("panel/update", controllers.PanelUpdate)
				// 面板删除
				admin.DELETE("panel/delete", controllers.PanelDelete)
				// 面板测试连接
				admin.POST("panel/test/connect", controllers.PanelTestConnect)
				// 面板绑定变量
				admin.PUT("panel/binding/update", controllers.PanelBindingUpdate)
				// 面板解除所有变量绑定
				admin.PUT("panel/unbind/update", controllers.PanelUnbindUpdate)
				// 面板批量更新Token
				admin.PUT("panel/token/update", controllers.PanelTokenUpdate)
			}

			// 变量管理
			{
				// 变量面板分页查询
				admin.GET("env/division/data", controllers.EnvDivisionData)
				// 获取变量简易数据
				admin.GET("env/division/data/simple", controllers.EnvDivisionDataSimple)
				// 变量新增
				admin.POST("env/add", controllers.EnvAdd)
				// 变量修改
				admin.PUT("env/update", controllers.EnvUpdate)
				// 变量删除
				admin.DELETE("env/delete", controllers.EnvDelete)
				// 手动刷新首页缓存数据
				admin.PUT("env/refresh", controllers.UserRefresh)
			}

			// 消息推送管理
			{
				// 用户WxPusher分页查询（仅返回已绑定WxPusher的用户信息）
				admin.GET("message/division/data", controllers.MessageDivisionData)
				// 管理员消息群发
				admin.POST("message/send", controllers.MessageSend)
				// 管理员全体消息发送
				admin.POST("message/send/all", controllers.MessageSendAll)
			}

			// 容器管理
			{
				// 迁移
				admin.POST("container/transfer", controllers.ContainerTransfer)
				// 复制
				admin.POST("container/copy", controllers.ContainerCopy)
				// 备份
				admin.POST("container/backup", controllers.ContainerBackup)
				// 恢复
				admin.POST("container/restore", controllers.ContainerRestore)
				// 备份数据下载
				admin.GET("container/backup/download", controllers.ContainerBackupDownload)
				// config.sh 同步
				admin.POST("container/synchronization", controllers.ContainerSynchronization)
				// 任务错误记录查询
				admin.GET("container/error/content", controllers.ContainerErrorContent)
				// 定时备份面板变量
				admin.POST("container/cron/backup", controllers.ContainerCronBackup)
			}

			// 定时任务
			{
				// 分页查询
				//admin.GET("cron/task/division/data", controllers.CronTaskDivisionData)
				// 创建任务
				//admin.POST("cron/task/add", controllers.CronTaskAdd)
				// 修改任务
				//admin.PUT("cron/task/update")
				// 删除任务
				//admin.DELETE("cron/task/delete")
			}

			// 插件管理
			{
				// 读取插件（上传执行插件、定时运行插件）
				admin.GET("plugin/data", controllers.PluginData)
				// 刷新定时插件
				admin.PUT("plugin/cron/refresh", controllers.PluginCronRefresh)
				// 上传插件（上传执行插件、定时运行插件）
				admin.POST("plugin/upload", controllers.PluginUpload)
				// 删除插件
				admin.DELETE("plugin/delete", controllers.PluginDelete)
				// 下载远程插件库插件
				admin.POST("plugin/remote/download", controllers.PluginRemoteDownload)
			}

			// 用户管理
			{
				// 分页查询
				admin.GET("user/division/data", controllers.UserDivisionData)
				// 筛选查询
				admin.GET("user/search", controllers.UserSearch)
				// 修改
				admin.PUT("user/information/update", controllers.UserInformationUpdate)
				// 删除
				admin.DELETE("user/information/delete", controllers.UserInformationDelete)
			}

			// 用户变量管理
			{
				// 筛选查询
				admin.GET("user/env/data/search", controllers.UserEnvDataSearch)
				// 修改用户变量
				admin.PUT("user/env/data/update", controllers.UserEnvDataUpdate)
				// 删除用户变量
				admin.DELETE("user/env/data/delete", controllers.UserEnvDataDelete)
			}

			// 卡密管理
			{
				// 分页查询
				admin.GET("cd-key/division/data", controllers.CDKEYDivisionCDKData)
				// 筛选查询(integral：积分、vip：会员)
				admin.GET("cd-key/search", controllers.CDKEYSearch)
				// 标识查询
				admin.GET("cd-key/remarks/search", controllers.CDKEYRemarksSearch)
				// 生成
				admin.POST("cd-key/add", controllers.CDKEYAdd)
				// 下载文件
				admin.GET("cd-key/data/download", controllers.CDKEYDataDownload)
				// 修改
				admin.PUT("cd-key/batch/update", controllers.CDKEYBatchUpdate)
				// 删除
				admin.DELETE("cd-key/delete", controllers.CDKEYDelete)
				// 指定用户充值
				admin.POST("cd-key/user/recharge", controllers.CDKEYUserRechargeIntegral)
			}

			// 充值数据
			{
				// 分页查询
				admin.GET("recharge/division/data", controllers.RechargeDivisionData)
				// 筛选查询
				admin.GET("recharge/search", controllers.RechargeSearch)
			}

			// 上传记录
			{
				// 分页查询
				admin.GET("record/division/data", controllers.RecordDivisionData)
				// 筛选查询
				admin.GET("record/search", controllers.RecordSearch)
			}

			// 网站设置
			{
				// 设置修改
				admin.PUT("set/settings", controllers.SaveSettings)
			}
		}
	}

	return r
}

// loadTemplate 加载模板文件
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range bindata.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := bindata.Asset(name)
		if err != nil {
			continue
		}
		name := strings.Replace(name, "assets/", "", 1)
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
