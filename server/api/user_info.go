package api

import (
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/model/request"
	"Kama-Chat/service"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserInfo 用户信息
var UserInfo = &UserInfoController{}

// UserInfoController 控制器
type UserInfoController struct {
	// userInfoSrv 用户信息服务
	userInfoSrv *service.UserInfoService
}

// Register 注册
func (uic *UserInfoController) Register(c *gin.Context) {
	req := &request.RegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.Register(req)
	response.JsonBack(c, message, ret, userInfo)
}

// Login 登录
func (uic *UserInfoController) Login(c *gin.Context) {
	req := &request.LoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.Login(req)
	response.JsonBack(c, message, ret, userInfo)
}

// UpdateUserInfo 修改用户信息
func (uic *UserInfoController) UpdateUserInfo(c *gin.Context) {
	req := &request.UpdateUserInfoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.UpdateUserInfo(req)
	response.JsonBack(c, message, ret, nil)
}

// GetUserInfoList 获取用户列表
func (uic *UserInfoController) GetUserInfoList(c *gin.Context) {
	req := &request.GetUserInfoListRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userList, ret := uic.userInfoSrv.GetUserInfoList(req)
	response.JsonBack(c, message, ret, userList)
}

// AbleUsers 启用用户
func (uic *UserInfoController) AbleUsers(c *gin.Context) {
	req := &request.AbleUsersRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.AbleUsers(req)
	response.JsonBack(c, message, ret, nil)
}

// DisableUsers 禁用用户
func (uic *UserInfoController) DisableUsers(c *gin.Context) {
	req := &request.AbleUsersRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.DisableUsers(req)
	response.JsonBack(c, message, ret, nil)
}

// GetUserInfo 获取用户信息
func (uic *UserInfoController) GetUserInfo(c *gin.Context) {
	req := &request.GetUserInfoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.GetUserInfo(req)
	response.JsonBack(c, message, ret, userInfo)
}

// DeleteUsers 删除用户
func (uic *UserInfoController) DeleteUsers(c *gin.Context) {
	req := &request.AbleUsersRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.DeleteUsers(req)
	response.JsonBack(c, message, ret, nil)
}

// SetAdmin 设置管理员
func (uic *UserInfoController) SetAdmin(c *gin.Context) {
	req := &request.AbleUsersRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.SetAdmin(req)
	response.JsonBack(c, message, ret, nil)
}

// SmsLogin 验证码登录
func (uic *UserInfoController) SmsLogin(c *gin.Context) {
	req := &request.SmsLoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.SmsLogin(req)
	response.JsonBack(c, message, ret, userInfo)
}

// SendSmsCode 发送短信验证码
func (uic *UserInfoController) SendSmsCode(c *gin.Context) {
	req := &request.SendSmsCodeRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := uic.userInfoSrv.SendSmsCode(req)
	response.JsonBack(c, message, ret, nil)
}
