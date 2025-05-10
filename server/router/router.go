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
	userGp := Router.Group("/user")
	{
		userGp.POST("/update_user", api.UserInfo.UpdateUserInfo)
		userGp.POST("/get_user_list", api.UserInfo.GetUserInfoList)
		userGp.POST("/able_users", api.UserInfo.AbleUsers)
		userGp.POST("/get_user_info", api.UserInfo.GetUserInfo)
		userGp.POST("/disable_users", api.UserInfo.DisableUsers)
		userGp.POST("/delete_users", api.UserInfo.DeleteUsers)
		userGp.POST("/set_admin", api.UserInfo.SetAdmin)
		userGp.POST("/send_sms_code", api.UserInfo.SendSmsCode)
		userGp.POST("/smsLogin", api.UserInfo.SmsLogin)
	}

	groupGp := Router.Group("/group")
	{
		groupGp.POST("/create_group", api.GroupInfo.CreateGroup)
		groupGp.POST("/load_my_group", api.GroupInfo.LoadMyGroup)
		groupGp.POST("/check_group_add_mode", api.GroupInfo.CheckGroupAddMode)
		groupGp.POST("/enter_group_directly", api.GroupInfo.EnterGroupDirectly)
		groupGp.POST("/leave_group", api.GroupInfo.LeaveGroup)
		//groupGp.POST("/dismissGroup", api.DismissGroup)
		//groupGp.POST("/getGroupInfo", api.GetGroupInfo)
		//groupGp.POST("/getGroupInfoList", api.GetGroupInfoList)
		//groupGp.POST("/deleteGroups", api.DeleteGroups)
		//groupGp.POST("/setGroupsStatus", api.SetGroupsStatus)
		//groupGp.POST("/updateGroupInfo", api.UpdateGroupInfo)
		//groupGp.POST("/getGroupMemberList", api.GetGroupMemberList)
		//groupGp.POST("/removeGroupMembers", api.RemoveGroupMembers)
	}
	//
	//sessionGp := Router.Group("/session")
	//{
	//	sessionGp.POST("/openSession", api.OpenSession)
	//	sessionGp.POST("/getUserSessionList", api.GetUserSessionList)
	//	sessionGp.POST("/getGroupSessionList", api.GetGroupSessionList)
	//	sessionGp.POST("/deleteSession", api.DeleteSession)
	//	sessionGp.POST("/checkOpenSessionAllowed", api.CheckOpenSessionAllowed)
	//}
	//
	//contactGp := Router.Group("/contact")
	//{
	//	contactGp.POST("/getUserList", api.GetUserList)
	//	contactGp.POST("/loadMyJoinedGroup", api.LoadMyJoinedGroup)
	//	contactGp.POST("/getContactInfo", api.GetContactInfo)
	//	contactGp.POST("/deleteContact", api.DeleteContact)
	//	contactGp.POST("/applyContact", api.ApplyContact)
	//	contactGp.POST("/getNewContactList", api.GetNewContactList)
	//	contactGp.POST("/passContactApply", api.PassContactApply)
	//	contactGp.POST("/blackContact", api.BlackContact)
	//	contactGp.POST("/cancelBlackContact", api.CancelBlackContact)
	//	contactGp.POST("/getAddGroupList", api.GetAddGroupList)
	//	contactGp.POST("/refuseContactApply", api.RefuseContactApply)
	//	contactGp.POST("/blackApply", api.BlackApply)
	//}
	//
	//messageGp := Router.Group("/message")
	//{
	//	messageGp.POST("/getMessageList", api.GetMessageList)
	//	messageGp.POST("/getGroupMessageList", api.GetGroupMessageList)
	//	messageGp.POST("/uploadAvatar", api.UploadAvatar)
	//	messageGp.POST("/uploadFile", api.UploadFile)
	//}

	//Router.POST("/chatroom/getCurContactListInChatRoom", api.GetCurContactListInChatRoom)
	//Router.GET("/wss", api.WsLogin)
	//Router.POST("/user/wsLogout", api.WsLogout)

}
