package api

import (
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/model/request"
	"Kama-Chat/service"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var UserContact = &UserContactController{}

type UserContactController struct {
	userContactSrv *service.UserContactService
}

// GetUserContactList 获取联系人列表
func (ucc *UserContactController) GetUserContactList(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
	}
	message, userList, ret := ucc.userContactSrv.GetUserList(req)
	response.JsonBack(c, message, ret, userList)
}

// LoadMyJoinedGroup 获取我加入的群聊
func (ucc *UserContactController) LoadMyJoinedGroup(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupList, ret := ucc.userContactSrv.LoadMyJoinedGroup(req)
	response.JsonBack(c, message, ret, groupList)
}

// GetContactInfo 获取联系人信息
func (ucc *UserContactController) GetContactInfo(c *gin.Context) {
	req := &request.GetContactInfoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	log.Println(req)
	message, contactInfo, ret := ucc.userContactSrv.GetContactInfo(req)
	response.JsonBack(c, message, ret, contactInfo)
}
