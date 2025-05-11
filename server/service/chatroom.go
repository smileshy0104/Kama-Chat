package service

import (
	"Kama-Chat/model/request"
	"Kama-Chat/model/respond"
	"github.com/gin-gonic/gin"
)

type ChatRoomService struct {
	Ctx *gin.Context
}

type chatRoomKey struct {
	ownerId   string
	contactId string
}

// map 类型是 {string, string}: []string该怎么写
var chatRooms = make(map[chatRoomKey][]string)

// GetCurContactListInChatRoom 获取当前聊天室联系人列表
func (crs *ChatRoomService) GetCurContactListInChatRoom(req *request.GetCurContactListInChatRoomRequest) (string, []respond.GetCurContactListInChatRoomRespond, int) {
	var rspList []respond.GetCurContactListInChatRoomRespond
	for _, contactId := range chatRooms[chatRoomKey{req.OwnerId, req.ContactId}] {
		rspList = append(rspList, respond.GetCurContactListInChatRoomRespond{
			ContactId: contactId,
		})
	}
	return "获取聊天室联系人列表成功", rspList, 0
}
