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

// DeleteContact 删除联系人
func (ucc *UserContactController) DeleteContact(c *gin.Context) {
	req := &request.DeleteContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.DeleteContact(req)
	response.JsonBack(c, message, ret, nil)
}

// ApplyContact 申请添加联系人
func (ucc *UserContactController) ApplyContact(c *gin.Context) {
	req := &request.ApplyContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.ApplyContact(req)
	response.JsonBack(c, message, ret, nil)
}

// GetNewContactList 获取新的联系人申请列表
func (ucc *UserContactController) GetNewContactList(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, data, ret := ucc.userContactSrv.GetNewContactList(req)
	response.JsonBack(c, message, ret, data)
}

// PassContactApply 通过联系人申请
func (ucc *UserContactController) PassContactApply(c *gin.Context) {
	req := &request.PassContactApplyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.PassContactApply(req)
	response.JsonBack(c, message, ret, nil)
}

// RefuseContactApply 拒绝联系人申请
func (ucc *UserContactController) RefuseContactApply(c *gin.Context) {
	req := &request.PassContactApplyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.RefuseContactApply(req)
	response.JsonBack(c, message, ret, nil)
}

// BlackContact 拉黑联系人
func (ucc *UserContactController) BlackContact(c *gin.Context) {
	req := &request.BlackContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.BlackContact(req)
	response.JsonBack(c, message, ret, nil)
}

// CancelBlackContact 解除拉黑联系人
func (ucc *UserContactController) CancelBlackContact(c *gin.Context) {
	req := &request.BlackContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.CancelBlackContact(req)
	response.JsonBack(c, message, ret, nil)
}

// GetAddGroupList 获取新的群聊申请列表
func (ucc *UserContactController) GetAddGroupList(c *gin.Context) {
	req := &request.AddGroupListRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, data, ret := ucc.userContactSrv.GetAddGroupList(req)
	response.JsonBack(c, message, ret, data)
}

// BlackApply 拉黑申请
func (ucc *UserContactController) BlackApply(c *gin.Context) {
	req := &request.BlackApplyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := ucc.userContactSrv.BlackApply(req)
	response.JsonBack(c, message, ret, nil)
}
