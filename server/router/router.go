package router

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/ssl"
	"Kama-Chat/middleware"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"*"}
	//corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	//corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	Router.Use(middleware.Cors())
	Router.Use(ssl.TlsHandler(global.CONFIG.MainConfig.Host, global.CONFIG.MainConfig.Port))
	Router.Static("/static/avatars", global.CONFIG.StaticSrcConfig.StaticAvatarPath)
	Router.Static("/static/files", global.CONFIG.StaticSrcConfig.StaticFilePath)
	Router.POST("/login", v1.Login)
	Router.POST("/register", v1.Register)
	Router.POST("/user/updateUserInfo", v1.UpdateUserInfo)
	Router.POST("/user/getUserInfoList", v1.GetUserInfoList)
	Router.POST("/user/ableUsers", v1.AbleUsers)
	Router.POST("/user/getUserInfo", v1.GetUserInfo)
	Router.POST("/user/disableUsers", v1.DisableUsers)
	Router.POST("/user/deleteUsers", v1.DeleteUsers)
	Router.POST("/user/setAdmin", v1.SetAdmin)
	Router.POST("/user/sendSmsCode", v1.SendSmsCode)
	Router.POST("/user/smsLogin", v1.SmsLogin)
	Router.POST("/user/wsLogout", v1.WsLogout)
	Router.POST("/group/createGroup", v1.CreateGroup)
	Router.POST("/group/loadMyGroup", v1.LoadMyGroup)
	Router.POST("/group/checkGroupAddMode", v1.CheckGroupAddMode)
	Router.POST("/group/enterGroupDirectly", v1.EnterGroupDirectly)
	Router.POST("/group/leaveGroup", v1.LeaveGroup)
	Router.POST("/group/dismissGroup", v1.DismissGroup)
	Router.POST("/group/getGroupInfo", v1.GetGroupInfo)
	Router.POST("/group/getGroupInfoList", v1.GetGroupInfoList)
	Router.POST("/group/deleteGroups", v1.DeleteGroups)
	Router.POST("/group/setGroupsStatus", v1.SetGroupsStatus)
	Router.POST("/group/updateGroupInfo", v1.UpdateGroupInfo)
	Router.POST("/group/getGroupMemberList", v1.GetGroupMemberList)
	Router.POST("/group/removeGroupMembers", v1.RemoveGroupMembers)
	Router.POST("/session/openSession", v1.OpenSession)
	Router.POST("/session/getUserSessionList", v1.GetUserSessionList)
	Router.POST("/session/getGroupSessionList", v1.GetGroupSessionList)
	Router.POST("/session/deleteSession", v1.DeleteSession)
	Router.POST("/session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)
	Router.POST("/contact/getUserList", v1.GetUserList)
	Router.POST("/contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
	Router.POST("/contact/getContactInfo", v1.GetContactInfo)
	Router.POST("/contact/deleteContact", v1.DeleteContact)
	Router.POST("/contact/applyContact", v1.ApplyContact)
	Router.POST("/contact/getNewContactList", v1.GetNewContactList)
	Router.POST("/contact/passContactApply", v1.PassContactApply)
	Router.POST("/contact/blackContact", v1.BlackContact)
	Router.POST("/contact/cancelBlackContact", v1.CancelBlackContact)
	Router.POST("/contact/getAddGroupList", v1.GetAddGroupList)
	Router.POST("/contact/refuseContactApply", v1.RefuseContactApply)
	Router.POST("/contact/blackApply", v1.BlackApply)
	Router.POST("/message/getMessageList", v1.GetMessageList)
	Router.POST("/message/getGroupMessageList", v1.GetGroupMessageList)
	Router.POST("/message/uploadAvatar", v1.UploadAvatar)
	Router.POST("/message/uploadFile", v1.UploadFile)
	Router.POST("/chatroom/getCurContactListInChatRoom", v1.GetCurContactListInChatRoom)
	Router.GET("/wss", v1.WsLogin)

}
