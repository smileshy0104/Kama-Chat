package service

import (
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	myredis "Kama-Chat/lib/redis"
	"Kama-Chat/model"
	"Kama-Chat/model/request"
	"Kama-Chat/model/respond"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/enum"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

type UserContactService struct {
	Ctx *gin.Context
}

// GetUserList 获取用户列表
// 关于用户被禁用的问题，这里查到的是所有联系人，如果被禁用或被拉黑会以弹窗的形式提醒，无法打开会话框；如果被删除，是搜索不到该联系人的。
func (ucs *UserContactService) GetUserList(req *request.OwnlistRequest) (string, []respond.MyUserListRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("contact_user_list_" + req.OwnerId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// dao
			var contactList []model.UserContact
			// 没有被删除
			if res := dao.GormDB.Order("created_at DESC").Where("user_id = ? AND status != 4", req.OwnerId).Find(&contactList); res.Error != nil {
				// 不存在不是业务问题，用Info，return 0
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					message := "目前不存在联系人"
					zlog.Info(message)
					return message, nil, 0
				} else {
					zlog.Error(res.Error.Error())
					return constants.SYSTEM_ERROR, nil, -1
				}
			}
			// dto
			var userListRsp []respond.MyUserListRespond
			for _, contact := range contactList {
				// 联系人中是用户的
				if contact.ContactType == enum.USER {
					// 获取用户信息
					var user model.UserInfo
					if res := dao.GormDB.First(&user, "uuid = ?", contact.ContactId); res.Error != nil {
						// 肯定是存在的，不可能无缘无故删掉，目前不用加notfound的判断
						zlog.Error(res.Error.Error())
						return constants.SYSTEM_ERROR, nil, -1
					}
					userListRsp = append(userListRsp, respond.MyUserListRespond{
						UserId:   user.Uuid,
						UserName: user.Nickname,
						Avatar:   user.Avatar,
					})
				}
			}
			rspString, err := json.Marshal(userListRsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("contact_user_list_"+req.OwnerId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取用户列表成功", userListRsp, 0
		} else {
			zlog.Error(err.Error())
		}
	}
	var rsp []respond.MyUserListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取用户列表成功", rsp, 0
}

// LoadMyJoinedGroup 获取我加入的群聊
func (ucs *UserContactService) LoadMyJoinedGroup(req *request.OwnlistRequest) (string, []respond.LoadMyJoinedGroupRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("my_joined_group_list_" + req.OwnerId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var contactList []model.UserContact
			// 没有退群，也没有被踢出群聊
			if res := dao.GormDB.Order("created_at DESC").Where("user_id = ? AND status != 6 AND status != 7", req.OwnerId).Find(&contactList); res.Error != nil {
				// 不存在不是业务问题，用Info，return 0
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					message := "目前不存在加入的群聊"
					zlog.Info(message)
					return message, nil, 0
				} else {
					zlog.Error(res.Error.Error())
					return constants.SYSTEM_ERROR, nil, -1
				}
			}
			var groupList []model.GroupInfo
			for _, contact := range contactList {
				if contact.ContactId[0] == 'G' {
					// 获取群聊信息
					var group model.GroupInfo
					if res := dao.GormDB.First(&group, "uuid = ?", contact.ContactId); res.Error != nil {
						zlog.Error(res.Error.Error())
						return constants.SYSTEM_ERROR, nil, -1
					}
					// 群没被删除，同时群主不是自己
					// 群主删除或admin删除群聊，status为7，即被踢出群聊，所以不用判断群是否被删除，删除了到不了这步
					if group.OwnerId != req.OwnerId {
						groupList = append(groupList, group)
					}
				}
			}
			var groupListRsp []respond.LoadMyJoinedGroupRespond
			for _, group := range groupList {
				groupListRsp = append(groupListRsp, respond.LoadMyJoinedGroupRespond{
					GroupId:   group.Uuid,
					GroupName: group.Name,
					Avatar:    group.Avatar,
				})
			}
			rspString, err := json.Marshal(groupListRsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("my_joined_group_list_"+req.OwnerId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取加入群成功", groupListRsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp []respond.LoadMyJoinedGroupRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取加入群成功", rsp, 0
}

// GetContactInfo 获取联系人信息
// 调用这个接口的前提是该联系人没有处在删除或被删除，或者该用户还在群聊中
// redis todo
func (ucs *UserContactService) GetContactInfo(req *request.GetContactInfoRequest) (string, respond.GetContactInfoRespond, int) {
	if req.ContactId[0] == 'G' {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, "uuid = ?", req.ContactId); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, respond.GetContactInfoRespond{}, -1
		}
		// 没被禁用
		if group.Status != enum.DISABLE {
			return "获取联系人信息成功", respond.GetContactInfoRespond{
				ContactId:        group.Uuid,
				ContactName:      group.Name,
				ContactAvatar:    group.Avatar,
				ContactNotice:    group.Notice,
				ContactAddMode:   group.AddMode,
				ContactMembers:   group.Members,
				ContactMemberCnt: group.MemberCnt,
				ContactOwnerId:   group.OwnerId,
			}, 0
		} else {
			zlog.Error("该群聊处于禁用状态")
			return "该群聊处于禁用状态", respond.GetContactInfoRespond{}, -2
		}
	} else {
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", req.ContactId); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, respond.GetContactInfoRespond{}, -1
		}
		log.Println(user)
		if user.Status != enum.DISABLE {
			return "获取联系人信息成功", respond.GetContactInfoRespond{
				ContactId:        user.Uuid,
				ContactName:      user.Nickname,
				ContactAvatar:    user.Avatar,
				ContactBirthday:  user.Birthday,
				ContactEmail:     user.Email,
				ContactPhone:     user.Telephone,
				ContactGender:    user.Gender,
				ContactSignature: user.Signature,
			}, 0
		} else {
			zlog.Info("该用户处于禁用状态")
			return "该用户处于禁用状态", respond.GetContactInfoRespond{}, -2
		}
	}
}
