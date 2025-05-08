package route

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/ssl"
	"Kama-Chat/middleware"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	r = gin.Default()
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"*"}
	//corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	//corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(middleware.Cors())
	r.Use(ssl.TlsHandler(global.CONFIG.MainConfig.Host, global.CONFIG.MainConfig.Port))
	r.Static("/static/avatars", global.CONFIG.StaticSrcConfig.StaticAvatarPath)
	r.Static("/static/files", global.CONFIG.StaticSrcConfig.StaticFilePath)
	r.POST("/login", v1.Login)
	r.POST("/register", v1.Register)
	r.POST("/user/updateUserInfo", v1.UpdateUserInfo)
	r.POST("/user/getUserInfoList", v1.GetUserInfoList)
	r.POST("/user/ableUsers", v1.AbleUsers)
	r.POST("/user/getUserInfo", v1.GetUserInfo)
	r.POST("/user/disableUsers", v1.DisableUsers)
	r.POST("/user/deleteUsers", v1.DeleteUsers)
	r.POST("/user/setAdmin", v1.SetAdmin)
	r.POST("/user/sendSmsCode", v1.SendSmsCode)
	r.POST("/user/smsLogin", v1.SmsLogin)
	r.POST("/user/wsLogout", v1.WsLogout)
	r.POST("/group/createGroup", v1.CreateGroup)
	r.POST("/group/loadMyGroup", v1.LoadMyGroup)
	r.POST("/group/checkGroupAddMode", v1.CheckGroupAddMode)
	r.POST("/group/enterGroupDirectly", v1.EnterGroupDirectly)
	r.POST("/group/leaveGroup", v1.LeaveGroup)
	r.POST("/group/dismissGroup", v1.DismissGroup)
	r.POST("/group/getGroupInfo", v1.GetGroupInfo)
	r.POST("/group/getGroupInfoList", v1.GetGroupInfoList)
	r.POST("/group/deleteGroups", v1.DeleteGroups)
	r.POST("/group/setGroupsStatus", v1.SetGroupsStatus)
	r.POST("/group/updateGroupInfo", v1.UpdateGroupInfo)
	r.POST("/group/getGroupMemberList", v1.GetGroupMemberList)
	r.POST("/group/removeGroupMembers", v1.RemoveGroupMembers)
	r.POST("/session/openSession", v1.OpenSession)
	r.POST("/session/getUserSessionList", v1.GetUserSessionList)
	r.POST("/session/getGroupSessionList", v1.GetGroupSessionList)
	r.POST("/session/deleteSession", v1.DeleteSession)
	r.POST("/session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)
	r.POST("/contact/getUserList", v1.GetUserList)
	r.POST("/contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
	r.POST("/contact/getContactInfo", v1.GetContactInfo)
	r.POST("/contact/deleteContact", v1.DeleteContact)
	r.POST("/contact/applyContact", v1.ApplyContact)
	r.POST("/contact/getNewContactList", v1.GetNewContactList)
	r.POST("/contact/passContactApply", v1.PassContactApply)
	r.POST("/contact/blackContact", v1.BlackContact)
	r.POST("/contact/cancelBlackContact", v1.CancelBlackContact)
	r.POST("/contact/getAddGroupList", v1.GetAddGroupList)
	r.POST("/contact/refuseContactApply", v1.RefuseContactApply)
	r.POST("/contact/blackApply", v1.BlackApply)
	r.POST("/message/getMessageList", v1.GetMessageList)
	r.POST("/message/getGroupMessageList", v1.GetGroupMessageList)
	r.POST("/message/uploadAvatar", v1.UploadAvatar)
	r.POST("/message/uploadFile", v1.UploadFile)
	r.POST("/chatroom/getCurContactListInChatRoom", v1.GetCurContactListInChatRoom)
	r.GET("/wss", v1.WsLogin)

}
