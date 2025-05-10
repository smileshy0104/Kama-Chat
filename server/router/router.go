package router

import (
	"Kama-Chat/api"
	"Kama-Chat/global"
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
	//Router.Use(ssl.TlsHandler(global.CONFIG.MainConfig.Host, global.CONFIG.MainConfig.Port))
	Router.Static("/static/avatars", global.CONFIG.StaticSrcConfig.StaticAvatarPath)
	Router.Static("/static/files", global.CONFIG.StaticSrcConfig.StaticFilePath)
	Router.POST("/register", api.UserInfo.Register)
	Router.POST("/login", api.UserInfo.Login)
	Router.POST("/user/update_user", api.UserInfo.UpdateUserInfo)
	Router.POST("/user/get_user_list", api.UserInfo.GetUserInfoList)
	Router.POST("/user/able_users", api.UserInfo.AbleUsers)
	Router.POST("/user/get_user_info", api.UserInfo.GetUserInfo)
	Router.POST("/user/disable_users", api.UserInfo.DisableUsers)
	Router.POST("/user/delete_users", api.UserInfo.DeleteUsers)
	Router.POST("/user/set_admin", api.UserInfo.SetAdmin)
	Router.POST("/user/sendSmsCode", api.UserInfo.SendSmsCode)
	Router.POST("/user/smsLogin", api.UserInfo.SmsLogin)
	//Router.POST("/user/wsLogout", api.WsLogout)
	//Router.POST("/group/createGroup", api.CreateGroup)
	//Router.POST("/group/loadMyGroup", api.LoadMyGroup)
	//Router.POST("/group/checkGroupAddMode", api.CheckGroupAddMode)
	//Router.POST("/group/enterGroupDirectly", api.EnterGroupDirectly)
	//Router.POST("/group/leaveGroup", api.LeaveGroup)
	//Router.POST("/group/dismissGroup", api.DismissGroup)
	//Router.POST("/group/getGroupInfo", api.GetGroupInfo)
	//Router.POST("/group/getGroupInfoList", api.GetGroupInfoList)
	//Router.POST("/group/deleteGroups", api.DeleteGroups)
	//Router.POST("/group/setGroupsStatus", api.SetGroupsStatus)
	//Router.POST("/group/updateGroupInfo", api.UpdateGroupInfo)
	//Router.POST("/group/getGroupMemberList", api.GetGroupMemberList)
	//Router.POST("/group/removeGroupMembers", api.RemoveGroupMembers)
	//Router.POST("/session/openSession", api.OpenSession)
	//Router.POST("/session/getUserSessionList", api.GetUserSessionList)
	//Router.POST("/session/getGroupSessionList", api.GetGroupSessionList)
	//Router.POST("/session/deleteSession", api.DeleteSession)
	//Router.POST("/session/checkOpenSessionAllowed", api.CheckOpenSessionAllowed)
	//Router.POST("/contact/getUserList", api.GetUserList)
	//Router.POST("/contact/loadMyJoinedGroup", api.LoadMyJoinedGroup)
	//Router.POST("/contact/getContactInfo", api.GetContactInfo)
	//Router.POST("/contact/deleteContact", api.DeleteContact)
	//Router.POST("/contact/applyContact", api.ApplyContact)
	//Router.POST("/contact/getNewContactList", api.GetNewContactList)
	//Router.POST("/contact/passContactApply", api.PassContactApply)
	//Router.POST("/contact/blackContact", api.BlackContact)
	//Router.POST("/contact/cancelBlackContact", api.CancelBlackContact)
	//Router.POST("/contact/getAddGroupList", api.GetAddGroupList)
	//Router.POST("/contact/refuseContactApply", api.RefuseContactApply)
	//Router.POST("/contact/blackApply", api.BlackApply)
	//Router.POST("/message/getMessageList", api.GetMessageList)
	//Router.POST("/message/getGroupMessageList", api.GetGroupMessageList)
	//Router.POST("/message/uploadAvatar", api.UploadAvatar)
	//Router.POST("/message/uploadFile", api.UploadFile)
	//Router.POST("/chatroom/getCurContactListInChatRoom", api.GetCurContactListInChatRoom)
	//Router.GET("/wss", api.WsLogin)

}
