package api

import (
	"Kama-Chat/model/request"
	"Kama-Chat/service"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Message = &MessageController{}

type MessageController struct {
	messageSrv *service.MessageService
}

// GetMessageList 获取聊天记录
func (mc *MessageController) GetMessageList(c *gin.Context) {
	req := &request.GetMessageListRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := mc.messageSrv.GetMessageList(req)
	response.JsonBack(c, message, ret, rsp)
}

// GetGroupMessageList 获取群聊消息记录
func (mc *MessageController) GetGroupMessageList(c *gin.Context) {
	req := &request.GetGroupMessageListRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rsp, ret := mc.messageSrv.GetGroupMessageList(req)
	response.JsonBack(c, message, ret, rsp)
}

// UploadAvatar 上传头像
func (mc *MessageController) UploadAvatar(c *gin.Context) {
	message, ret := mc.messageSrv.UploadAvatar(c)
	response.JsonBack(c, message, ret, nil)
}

// UploadFile 上传头像
func (mc *MessageController) UploadFile(c *gin.Context) {
	message, ret := mc.messageSrv.UploadFile(c)
	response.JsonBack(c, message, ret, nil)
}
