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

var GroupInfo = &GroupInfoController{}

type GroupInfoController struct {
	groupInfoSrv *service.GroupInfoService
}

// CreateGroup 创建群聊
func (gic *GroupInfoController) CreateGroup(c *gin.Context) {
	req := &request.CreateGroupRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gic.groupInfoSrv.CreateGroup(req)
	response.JsonBack(c, message, ret, nil)
}

// LoadMyGroup 获取我创建的群聊
func (gic *GroupInfoController) LoadMyGroup(c *gin.Context) {
	req := &request.OwnlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupList, ret := gic.groupInfoSrv.LoadMyGroup(req)
	response.JsonBack(c, message, ret, groupList)
}

// CheckGroupAddMode 检查群聊加群方式
func (gic *GroupInfoController) CheckGroupAddMode(c *gin.Context) {
	req := &request.CheckGroupAddModeRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, addMode, ret := gic.groupInfoSrv.CheckGroupAddMode(req)
	response.JsonBack(c, message, ret, addMode)
}

// EnterGroupDirectly 直接进群
func (gic *GroupInfoController) EnterGroupDirectly(c *gin.Context) {
	req := &request.EnterGroupDirectlyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gic.groupInfoSrv.EnterGroupDirectly(req)
	response.JsonBack(c, message, ret, nil)
}

// LeaveGroup 退群
func (gic *GroupInfoController) LeaveGroup(c *gin.Context) {
	req := &request.LeaveGroupRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gic.groupInfoSrv.LeaveGroup(req)
	response.JsonBack(c, message, ret, nil)
}

// DismissGroup 解散群聊
func (gic *GroupInfoController) DismissGroup(c *gin.Context) {
	req := &request.DismissGroupRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gic.groupInfoSrv.DismissGroup(req)
	response.JsonBack(c, message, ret, nil)
}

// GetGroupInfo 获取群聊详情
func (gic *GroupInfoController) GetGroupInfo(c *gin.Context) {
	req := &request.GetGroupInfoRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupInfo, ret := gic.groupInfoSrv.GetGroupInfo(req)
	response.JsonBack(c, message, ret, groupInfo)
}

// GetGroupInfoList 获取群聊列表 - 管理员
func (gic *GroupInfoController) GetGroupInfoList(c *gin.Context) {
	message, groupList, ret := gic.groupInfoSrv.GetGroupInfoList()
	response.JsonBack(c, message, ret, groupList)
}

// DeleteGroups 删除列表中群聊 - 管理员
func (gic *GroupInfoController) DeleteGroups(c *gin.Context) {
	req := &request.DeleteGroupsRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gic.groupInfoSrv.DeleteGroups(req)
	response.JsonBack(c, message, ret, nil)
}
