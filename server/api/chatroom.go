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

var ChatRoom = &ChatRoomController{}

type ChatRoomController struct {
	chatRoomSrv *service.ChatRoomService
}

// GetCurContactListInChatRoom 获取当前聊天室联系人列表
func (crc *ChatRoomController) GetCurContactListInChatRoom(c *gin.Context) {
	req := &request.GetCurContactListInChatRoomRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, rspList, ret := crc.chatRoomSrv.GetCurContactListInChatRoom(req)
	response.JsonBack(c, message, ret, rspList)
}
