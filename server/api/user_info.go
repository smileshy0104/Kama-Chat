package api

import (
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/model/request"
	"Kama-Chat/service"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var UserInfo = &UserInfoController{}

type UserInfoController struct {
	userInfoSrv *service.UserInfoService
}

// Register 注册
func (uic *UserInfoController) Register(c *gin.Context) {
	registerReq := &request.RegisterRequest{}
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.Register(registerReq)
	response.JsonBack(c, message, ret, userInfo)
}

// Login 登录
func (uic *UserInfoController) Login(c *gin.Context) {
	loginReq := &request.LoginRequest{}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, userInfo, ret := uic.userInfoSrv.Login(loginReq)
	response.JsonBack(c, message, ret, userInfo)
}
