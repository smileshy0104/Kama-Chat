package service

import (
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	myredis "Kama-Chat/lib/redis"
	"Kama-Chat/lib/sms"
	"Kama-Chat/model"
	"Kama-Chat/model/request"
	"Kama-Chat/model/respond"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/enum"
	"Kama-Chat/utils/random"
	"Kama-Chat/utils/validate"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

// UserInfoService 用户信息服务
type UserInfoService struct {
	Ctx *gin.Context
}

// Register 注册
func (uis *UserInfoService) Register(req *request.RegisterRequest) (string, *respond.RegisterRespond, int) {
	key := "auth_code_" + req.Telephone
	// 验证码校验
	code, err := myredis.GetKey(key)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if code != req.SmsCode {
		message := "验证码不正确，请重试"
		zlog.Info(message)
		return message, nil, -2
	} else {
		// 删除已存在验证码
		if err := myredis.DelKeyIfExists(key); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	// 不用校验手机号，前端校验
	// 判断电话是否已经被注册过了
	message, ret := validate.CheckTelephoneExist(req.Telephone)
	if ret != 0 {
		return message, nil, ret
	}
	// 创建用户
	var newUser model.UserInfo
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.Telephone = req.Telephone
	newUser.Password = req.Password
	newUser.Nickname = req.Nickname
	newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = validate.CheckUserIsAdminOrNot(newUser)
	newUser.Status = enum.NORMAL

	// 创建用户
	res := dao.GormDB.Create(&newUser)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	// 注册成功，chat client建立
	//if err := chat.NewClientInit(c, newUser.Uuid); err != nil {
	//	return "", err
	//}

	// 注册成功，返回用户信息
	registerRsp := &respond.RegisterRespond{
		Uuid:      newUser.Uuid,
		Telephone: newUser.Telephone,
		Nickname:  newUser.Nickname,
		Email:     newUser.Email,
		Avatar:    newUser.Avatar,
		Gender:    newUser.Gender,
		Birthday:  newUser.Birthday,
		Signature: newUser.Signature,
		IsAdmin:   newUser.IsAdmin,
		Status:    newUser.Status,
	}
	year, month, day := newUser.CreatedAt.Date()
	registerRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "注册成功", registerRsp, 0
}

// Login 登录
func (uis *UserInfoService) Login(req *request.LoginRequest) (string, *respond.LoginRespond, int) {
	password := req.Password
	var user model.UserInfo
	// 获取用户信息
	res := dao.GormDB.First(&user, "telephone = ?", req.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	// 密码校验 TODO 密码目前未加密
	if user.Password != password {
		message := "密码不正确，请重试"
		zlog.Error(message)
		return message, nil, -2
	}
	// 登录成功，返回用户信息
	loginRsp := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Signature: user.Signature,
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	year, month, day := user.CreatedAt.Date()
	loginRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "登陆成功", loginRsp, 0
}

// UpdateUserInfo 修改用户信息
// 某用户修改了信息，可能会影响contact_user_list，不需要删除redis的contact_user_list，timeout之后会自己更新
// 但是需要更新redis的user_info，因为可能影响用户搜索
func (uis *UserInfoService) UpdateUserInfo(req *request.UpdateUserInfoRequest) (string, int) {
	var user model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.First(&user, "uuid = ?", req.Uuid); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Birthday != "" {
		user.Birthday = req.Birthday
	}
	if req.Signature != "" {
		user.Signature = req.Signature
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	// 保存到数据库
	if res := dao.GormDB.Save(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除redis
	if err := myredis.DelKeysWithPattern("user_info_" + req.Uuid); err != nil {
		zlog.Error(err.Error())
	}
	return "修改用户信息成功", 0
}

// GetUserInfoList 获取用户列表除了ownerId之外 - 管理员
// 管理员少，而且如果用户更改了，那么管理员会一直频繁删除redis，更新redis，比较麻烦，所以管理员暂时不使用redis缓存
func (uis *UserInfoService) GetUserInfoList(req *request.GetUserInfoListRequest) (string, []respond.GetUserListRespond, int) {
	// redis中没有数据，从数据库中获取
	var users []model.UserInfo
	// 获取所有的用户
	if res := dao.GormDB.Unscoped().Where("uuid != ?", req.OwnerId).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	// 获取用户信息
	var rsp []respond.GetUserListRespond
	for _, user := range users {
		rp := respond.GetUserListRespond{
			Uuid:      user.Uuid,
			Telephone: user.Telephone,
			Nickname:  user.Nickname,
			Status:    user.Status,
			IsAdmin:   user.IsAdmin,
		}
		if user.DeletedAt.Valid {
			rp.IsDeleted = true
		} else {
			rp.IsDeleted = false
		}
		rsp = append(rsp, rp)
	}
	return "获取用户列表成功", rsp, 0
}

// GetUserInfo 获取用户信息
func (uis *UserInfoService) GetUserInfo(req *request.GetUserInfoRequest) (string, *respond.GetUserInfoRespond, int) {
	zlog.Info(req.Uuid)
	// 从redis中获取user_info_（判断redis中是否存在对应的user_info）
	rspString, err := myredis.GetKeyNilIsErr("user_info_" + req.Uuid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info(err.Error())
			var user model.UserInfo
			// 获取用户信息
			if res := dao.GormDB.Where("uuid = ?", req.Uuid).Find(&user); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			// 获取用户信息
			rsp := respond.GetUserInfoRespond{
				Uuid:      user.Uuid,
				Telephone: user.Telephone,
				Nickname:  user.Nickname,
				Avatar:    user.Avatar,
				Birthday:  user.Birthday,
				Email:     user.Email,
				Gender:    user.Gender,
				Signature: user.Signature,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				IsAdmin:   user.IsAdmin,
				Status:    user.Status,
			}
			rspString, err := json.Marshal(rsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			// 存放user_info 到redis
			if err := myredis.SetKeyEx("user_info_"+req.Uuid, string(rspString), constants.REDIS_TIMEOUT*time.Minute); err != nil {
				zlog.Error(err.Error())
			}
			return "获取用户信息成功", &rsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp respond.GetUserInfoRespond
	// 解析json
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取用户信息成功", &rsp, 0
}

// AbleUsers 启用用户
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (uis *UserInfoService) AbleUsers(req *request.AbleUsersRequest) (string, int) {
	var users []model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", req.UuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 遍历更新状态
	for _, user := range users {
		user.Status = enum.NORMAL
		// 保存到数据库
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}
	// 删除所有"contact_user_list"开头的key
	if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
		zlog.Error(err.Error())
	}
	return "启用用户成功", 0
}

// DisableUsers 禁用用户——同时把联系人的会话删除
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (uis *UserInfoService) DisableUsers(req *request.AbleUsersRequest) (string, int) {
	var users []model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", req.UuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 遍历更新状态
	for _, user := range users {
		user.Status = enum.DISABLE
		// 保存到数据库
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		var sessionList []model.Session
		// 获取会话
		if res := dao.GormDB.Where("send_id = ? or receive_id = ?", user.Uuid, user.Uuid).Find(&sessionList); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		// 遍历更新会话
		for _, session := range sessionList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			session.DeletedAt = deletedAt
			// 保存会话
			if res := dao.GormDB.Save(&session); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
	}
	// 删除所有"contact_user_list"开头的key
	if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
		zlog.Error(err.Error())
	}
	return "禁用用户成功", 0
}

// DeleteUsers 删除用户
// 用户是否启用禁用需要实时更新contact_user_list状态，所以redis的contact_user_list需要删除
func (uis *UserInfoService) DeleteUsers(req *request.AbleUsersRequest) (string, int) {
	var users []model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", req.UuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 遍历更新状态
	for _, user := range users {
		user.DeletedAt.Valid = true
		user.DeletedAt.Time = time.Now()
		// 保存到数据库
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}

		// 删除会话
		var sessionList []model.Session
		if res := dao.GormDB.Where("send_id = ? or receive_id = ?", user.Uuid, user.Uuid).Find(&sessionList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		// 遍历更新会话
		for _, session := range sessionList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			session.DeletedAt = deletedAt
			// 保存会话
			if res := dao.GormDB.Save(&session); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

		// 删除联系人
		var contactList []model.UserContact
		if res := dao.GormDB.Where("user_id = ? or contact_id = ?", user.Uuid, user.Uuid).Find(&contactList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		// 遍历更新联系人
		for _, contact := range contactList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			contact.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&contact); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

		// 删除申请记录
		var applyList []model.ContactApply
		if res := dao.GormDB.Where("user_id = ? or contact_id = ?", user.Uuid, user.Uuid).Find(&applyList); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info(res.Error.Error())
			} else {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		// 遍历更新申请记录
		for _, apply := range applyList {
			var deletedAt gorm.DeletedAt
			deletedAt.Time = time.Now()
			deletedAt.Valid = true
			apply.DeletedAt = deletedAt
			if res := dao.GormDB.Save(&apply); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

	}
	// 删除所有"contact_user_list"开头的key
	if err := myredis.DelKeysWithPrefix("contact_user_list"); err != nil {
		zlog.Error(err.Error())
	}
	return "删除用户成功", 0
}

// SetAdmin 设置管理员
func (uis *UserInfoService) SetAdmin(req *request.AbleUsersRequest) (string, int) {
	var users []model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.Where("uuid = (?)", req.UuidList).Find(&users); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 遍历更新状态
	for _, user := range users {
		user.IsAdmin = req.IsAdmin
		if res := dao.GormDB.Save(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}
	return "设置管理员成功", 0
}

// SmsLogin 验证码登录
func (uis *UserInfoService) SmsLogin(req *request.SmsLoginRequest) (string, *respond.LoginRespond, int) {
	var user model.UserInfo
	// 查询用户
	res := dao.GormDB.First(&user, "telephone = ?", req.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 验证码校验
	key := "auth_code_" + req.Telephone
	code, err := myredis.GetKey(key)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if code != req.SmsCode {
		message := "验证码不正确，请重试"
		zlog.Info(message)
		return message, nil, -2
	} else {
		// 删除redis中的验证码
		if err := myredis.DelKeyIfExists(key); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}

	// 登录成功，返回响应
	loginRsp := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Signature: user.Signature,
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	year, month, day := user.CreatedAt.Date()
	loginRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "登陆成功", loginRsp, 0
}

// SendSmsCode 发送短信验证码 - 验证码登录
func (uis *UserInfoService) SendSmsCode(req *request.SendSmsCodeRequest) (string, int) {
	// 发送短信验证码
	return sms.VerificationCode(req.Telephone)
}
