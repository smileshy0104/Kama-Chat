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

var Session = &SessionController{}

type SessionController struct {
	sessionSrv *service.SessionService
}

// OpenSession 打开会话
func (sc *SessionController) OpenSession(c *gin.Context) {
	req := &request.OpenSessionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, sessionId, ret := sc.sessionSrv.OpenSession(req)
	response.JsonBack(c, message, ret, sessionId)
}

// GetUserSessionList 获取用户会话列表
func (sc *SessionController) GetUserSessionList(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, sessionList, ret := sc.sessionSrv.GetUserSessionList(req)
	response.JsonBack(c, message, ret, sessionList)
}

// GetGroupSessionList 获取群聊会话列表
func (sc *SessionController) GetGroupSessionList(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupList, ret := sc.sessionSrv.GetGroupSessionList(req)
	response.JsonBack(c, message, ret, groupList)
}

// DeleteSession 删除会话
func (sc *SessionController) DeleteSession(c *gin.Context) {
	req := &request.DeleteSessionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := sc.sessionSrv.DeleteSession(req)
	response.JsonBack(c, message, ret, nil)
}

// CheckOpenSessionAllowed 检查是否可以打开会话
func (sc *SessionController) CheckOpenSessionAllowed(c *gin.Context) {
	req := &request.CreateSessionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, res, ret := sc.sessionSrv.CheckOpenSessionAllowed(req)
	response.JsonBack(c, message, ret, res)
}
