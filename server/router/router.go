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

	// 用户相关
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

	// 群组相关
	groupGp := Router.Group("/group")
	{
		groupGp.POST("/create_group", api.GroupInfo.CreateGroup)
		groupGp.POST("/load_my_group", api.GroupInfo.LoadMyGroup)
		groupGp.POST("/check_group_add_mode", api.GroupInfo.CheckGroupAddMode)
		groupGp.POST("/enter_group_directly", api.GroupInfo.EnterGroupDirectly)
		groupGp.POST("/leave_group", api.GroupInfo.LeaveGroup)
		groupGp.POST("/dismiss_group", api.GroupInfo.DismissGroup)
		groupGp.POST("/get_group_info", api.GroupInfo.GetGroupInfo)
		groupGp.POST("/get_group_info_list", api.GroupInfo.GetGroupInfoList)
		groupGp.POST("/delete_groups", api.GroupInfo.DeleteGroups)
		groupGp.POST("/set_groups_status", api.GroupInfo.SetGroupsStatus)
		groupGp.POST("/update_group_info", api.GroupInfo.UpdateGroupInfo)
		groupGp.POST("/get_group_member_list", api.GroupInfo.GetGroupMemberList)
		groupGp.POST("/remove_group_members", api.GroupInfo.RemoveGroupMembers)
	}

	//sessionGp := Router.Group("/session")
	//{
	//	sessionGp.POST("/openSession", api.UserContact.OpenSession)
	//	sessionGp.POST("/getUserSessionList", api.UserContact.GetUserSessionList)
	//	sessionGp.POST("/getGroupSessionList", api.UserContact.GetGroupSessionList)
	//	sessionGp.POST("/deleteSession", api.UserContact.DeleteSession)
	//	sessionGp.POST("/checkOpenSessionAllowed", api.UserContact.CheckOpenSessionAllowed)
	//}
	//
	contactGp := Router.Group("/contact")
	{
		contactGp.POST("/getUserList", api.UserContact.GetUserContactList)
		contactGp.POST("/loadMyJoinedGroup", api.UserContact.LoadMyJoinedGroup)
		contactGp.POST("/getContactInfo", api.UserContact.GetContactInfo)
		contactGp.POST("/deleteContact", api.UserContact.DeleteContact)
		contactGp.POST("/applyContact", api.UserContact.ApplyContact)
		contactGp.POST("/getNewContactList", api.UserContact.GetNewContactList)
		contactGp.POST("/passContactApply", api.UserContact.PassContactApply)
		contactGp.POST("/blackContact", api.UserContact.BlackContact)
		contactGp.POST("/cancelBlackContact", api.UserContact.CancelBlackContact)
		//contactGp.POST("/getAddGroupList", api.UserContact.GetAddGroupList)
		//contactGp.POST("/refuseContactApply", api.UserContact.RefuseContactApply)
		//contactGp.POST("/blackApply", api.UserContact.BlackApply)
	}
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
