package router

import (
	"Kama-Chat/api"
	"Kama-Chat/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	Router.Use(cors.New(corsConfig))
	//Router.Use(middleware.Cors())
	//Router.Use(ssl.TlsHandler(global.CONFIG.MainConfig.Host, global.CONFIG.MainConfig.Port))
	Router.Static("/static/avatars", global.CONFIG.StaticSrcConfig.StaticAvatarPath)
	Router.Static("/static/files", global.CONFIG.StaticSrcConfig.StaticFilePath)
	Router.POST("/register", api.UserInfo.Register)
	Router.POST("/login", api.UserInfo.Login)

	// 用户相关
	userGp := Router.Group("/user")
	{
		userGp.POST("/update_user", api.UserInfo.UpdateUserInfo)
		userGp.POST("/get_user_info_list", api.UserInfo.GetUserInfoList)
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

	// 会话相关
	sessionGp := Router.Group("/session")
	{
		sessionGp.POST("/open_session", api.Session.OpenSession)
		sessionGp.POST("/get_user_session_list", api.Session.GetUserSessionList)
		sessionGp.POST("/get_group_session_list", api.Session.GetGroupSessionList)
		sessionGp.POST("/delete_session", api.Session.DeleteSession)
		sessionGp.POST("/check_open_session_allowed", api.Session.CheckOpenSessionAllowed)
	}

	// 联系人相关
	contactGp := Router.Group("/contact")
	{
		contactGp.POST("/get_user_contact_list", api.UserContact.GetUserContactList)
		contactGp.POST("/load_my_joined_group", api.UserContact.LoadMyJoinedGroup)
		contactGp.POST("/get_contact_info", api.UserContact.GetContactInfo)
		contactGp.POST("/delete_contact", api.UserContact.DeleteContact)
		contactGp.POST("/apply_contact", api.UserContact.ApplyContact)
		contactGp.POST("/get_new_contact_list", api.UserContact.GetNewContactList)
		contactGp.POST("/pass_contact_apply", api.UserContact.PassContactApply)
		contactGp.POST("/black_contact", api.UserContact.BlackContact)
		contactGp.POST("/cancel_black_contact", api.UserContact.CancelBlackContact)
		contactGp.POST("/get_add_group_list", api.UserContact.GetAddGroupList)
		contactGp.POST("/refuse_contact_apply", api.UserContact.RefuseContactApply)
		contactGp.POST("/black_apply", api.UserContact.BlackApply)
	}

	// 消息相关
	messageGp := Router.Group("/message")
	{
		messageGp.POST("/get_message_list", api.Message.GetMessageList)
		messageGp.POST("/get_group_message_list", api.Message.GetGroupMessageList)
		messageGp.POST("/upload_avatar", api.Message.UploadAvatar)
		messageGp.POST("/upload_file", api.Message.UploadFile)
	}

	// 聊天室相关
	chatRoomGp := Router.Group("/chatroom")
	{
		chatRoomGp.POST("/getCurContactListInChatRoom", api.ChatRoom.GetCurContactListInChatRoom)
	}

	Router.GET("/wss", api.Wss.WsLogin)
	Router.POST("/user/wsLogout", api.Wss.WsLogout)

}
