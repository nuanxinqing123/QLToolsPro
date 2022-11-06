// -*- coding: utf-8 -*-
// @Time    : 2022/8/18 10:59
// @Author  : Nuanxinqing
// @Email   : nuanxinqing@gmail.com
// @File    : user.go

package dao

import "QLToolsPro/server/model"

// GetUserRecord 获取User表记录数量
func GetUserRecord() int64 {
	return DB.Find(&model.User{}).RowsAffected
}

// CheckIsUser 校验是否属于注册用户
func CheckIsUser(userID interface{}) model.User {
	var user model.User
	DB.Where("user_id = ?", userID).First(&user)
	return user
}

// GetUserData 邮箱和用户名获取用户信息
func GetUserData(email, username string) (bool, model.User) {
	var user model.User
	// 通过邮箱和用户名查询用户信息
	DB.Where("email = ?", email).Or("username = ?", username).First(&user)

	// 判断是否已注册
	if user.Email != "" && user.Username != "" {
		// 存在
		return true, user
	} else {
		// 不存在
		return false, user
	}
}

// InsertUser 创建新用户
func InsertUser(user *model.User) (err error) {
	return DB.Create(&user).Error
}

// GetUserNameData 用户名获取用户信息
func GetUserNameData(username string) (bool, model.User) {
	var user model.User
	// 通过用户名查询用户信息
	DB.Where("username = ?", username).First(&user)

	// 判断是否已注册
	if user.Username != "" {
		// 存在
		return true, user
	} else {
		// 不存在
		return false, user
	}
}

// GetUserEmailData 邮箱获取用户信息
func GetUserEmailData(email string) (bool, model.User) {
	var user model.User
	// 通过邮箱和用户名查询用户信息
	DB.Where("email = ?", email).First(&user)

	// 判断是否已注册
	if user.Email != "" && user.Username != "" {
		// 存在
		return true, user
	} else {
		// 不存在
		return false, user
	}
}

// GetUserIDData 用户ID获取用户信息
func GetUserIDData(uid interface{}) model.User {
	var user model.User
	// 通过用户ID查询用户信息
	DB.Model(&model.User{}).Where("user_id = ?", uid).First(&user)
	return user
}

// GetUserAppOpenIDData 用户AppOpenID获取用户信息
func GetUserAppOpenIDData(openid string) model.User {
	var user model.User
	// 通过用户ID查询用户信息
	DB.Model(&model.User{}).Where("app_open_id = ?", openid).First(&user)
	return user
}

// UpdateUserLoginIP 更新用户登录IP地址
func UpdateUserLoginIP(ip string, uid interface{}) {
	var u model.User
	DB.Where("user_id = ?", uid).First(&u)
	u.LoginIP = ip
	DB.Save(&u)
}

// GetDivisionUserData 条件查询用户数据
func GetDivisionUserData(page, q int) []model.User {
	var user []model.User
	if page == 1 {
		// 获取第一页数据（25条）
		DB.Order("id desc").Limit(25).Offset(0).Find(&user)
	} else {
		// 获取第N页数据
		DB.Order("id desc").Limit(25).Offset((page - 1) * q).Find(&user)
	}
	return user
}

// GetUserNameSearch 用户名模糊查询
func GetUserNameSearch(s string) []model.User {
	var user []model.User
	DB.Where("username LIKE ?", "%"+s+"%").Find(&user)
	return user
}

// UserInformationUpdate 用户数据更新
func UserInformationUpdate(p *model.UpdateUserData) {
	user := new(model.User)
	DB.Where("id = ?", p.ID).First(&user)
	user.UserWxpusher = p.UserWxpusher
	user.IsAdmin = p.IsAdmin
	user.IsState = p.IsState
	DB.Save(&user)
}

// UserInformationDelete 删除User数据
func UserInformationDelete(p *model.DeleteUserData) {
	DB.Where("id = ?", p.ID).Delete(&model.User{})
}

// UpdateUserPwd 找回密码-修改密码
func UpdateUserPwd(uid string, pwd string) error {
	return DB.Model(&model.User{}).Where("user_id = ?", uid).Update("password", pwd).Error
}

// UpdateUserDataSave 用户数据保存
func UpdateUserDataSave(p model.User) {
	DB.Save(&p)
}

// GetUserWxpusherID 查询用户WxPusherID
func GetUserWxpusherID(uid interface{}) string {
	var user model.User
	// 通过用户ID查询用户信息
	DB.Model(&model.User{}).Where("user_id = ?", uid).First(&user)
	return user.UserWxpusher
}

// UpdateUserWxpusherID 更新用户WxPusherID
func UpdateUserWxpusherID(uid interface{}, wxUID string) error {
	return DB.Model(&model.User{}).Where("user_id = ?", uid).Update("user_wxpusher", wxUID).Error
}

// GetAllUserData 获取所有用户数据
func GetAllUserData() []model.User {
	var u []model.User
	DB.Find(&u)
	return u
}

// UpdateVIPState 修改用户VIP状态
func UpdateVIPState(u model.User) {
	DB.Save(&u)
}
