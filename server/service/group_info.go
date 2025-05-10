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
	"Kama-Chat/utils/random"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type GroupInfoService struct {
	Ctx *gin.Context
}

// CreateGroup 创建群聊
func (gis *GroupInfoService) CreateGroup(req *request.CreateGroupRequest) (string, int) {
	group := model.GroupInfo{
		Uuid:      fmt.Sprintf("G%s", random.GetNowAndLenRandomString(11)),
		Name:      req.Name,
		Notice:    req.Notice,
		OwnerId:   req.OwnerId,
		MemberCnt: 1,
		AddMode:   req.AddMode,
		Avatar:    req.Avatar,
		Status:    enum.NORMAL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var members []string
	members = append(members, req.OwnerId)
	var err error
	group.Members, err = json.Marshal(members)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if res := dao.GormDB.Create(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	// 添加联系人
	contact := model.UserContact{
		UserId:      req.OwnerId,
		ContactId:   group.Uuid,
		ContactType: enum.GROUP,
		Status:      enum.NORMAL,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
	if res := dao.GormDB.Create(&contact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if err := myredis.DelKeysWithPattern("contact_mygroup_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}

	return "创建成功", 0
}

// LoadMyGroup 获取我创建的群聊
func (gis *GroupInfoService) LoadMyGroup(req *request.OwnlistRequest) (string, []respond.LoadMyGroupRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("contact_mygroup_list_" + req.OwnerId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var groupList []model.GroupInfo
			if res := dao.GormDB.Order("created_at DESC").Where("owner_id = ?", req.OwnerId).Find(&groupList); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			var groupListRsp []respond.LoadMyGroupRespond
			for _, group := range groupList {
				groupListRsp = append(groupListRsp, respond.LoadMyGroupRespond{
					GroupId:   group.Uuid,
					GroupName: group.Name,
					Avatar:    group.Avatar,
				})
			}
			rspString, err := json.Marshal(groupListRsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("contact_mygroup_list_"+req.OwnerId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取成功", groupListRsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var groupListRsp []respond.LoadMyGroupRespond
	if err := json.Unmarshal([]byte(rspString), &groupListRsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取成功", groupListRsp, 0
}

// CheckGroupAddMode 检查群聊加群方式
func (gis *GroupInfoService) CheckGroupAddMode(req *request.CheckGroupAddModeRequest) (string, int8, int) {
	rspString, err := myredis.GetKeyNilIsErr("group_info_" + req.GroupId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var group model.GroupInfo
			if res := dao.GormDB.First(&group, "uuid = ?", req.GroupId); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1, -1
			}
			return "加群方式获取成功", group.AddMode, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1, -1
		}
	}
	var rsp respond.GetGroupInfoRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "加群方式获取成功", rsp.AddMode, 0
}

// EnterGroupDirectly 直接进群
// ownerId 是群聊id
func (gis *GroupInfoService) EnterGroupDirectly(req *request.EnterGroupDirectlyRequest) (string, int) {
	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", req.OwnerId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	var members []string
	if err := json.Unmarshal(group.Members, &members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	members = append(members, req.ContactId)
	if data, err := json.Marshal(members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	} else {
		group.Members = data
	}
	group.MemberCnt += 1
	if res := dao.GormDB.Save(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	newContact := model.UserContact{
		UserId:      req.ContactId,
		ContactId:   req.OwnerId,
		ContactType: enum.GROUP,   // 用户
		Status:      enum.NORMAL_, // 正常
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
	if res := dao.GormDB.Create(&newContact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	//if err := myredis.DelKeysWithPattern("group_info_" + contactId); err != nil {
	//	zlog.Error(err.Error())
	//}
	//if err := myredis.DelKeysWithPattern("groupmember_list_" + contactId); err != nil {
	//	zlog.Error(err.Error())
	//}
	if err := myredis.DelKeysWithPattern("group_session_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPattern("my_joined_group_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	//if err := myredis.DelKeysWithPattern("session_" + ownerId + "_" + contactId); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "进群成功", 0
}

// LeaveGroup 退群
func (gis *GroupInfoService) LeaveGroup(req *request.LeaveGroupRequest) (string, int) {
	// 从群聊中清除该用户
	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", req.GroupId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	var members []string
	if err := json.Unmarshal(group.Members, &members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for i, member := range members {
		if member == req.UserId {
			members = append(members[:i], members[i+1:]...)
			break
		}
	}
	if data, err := json.Marshal(members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	} else {
		group.Members = data
	}
	group.MemberCnt -= 1
	if res := dao.GormDB.Save(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除会话
	var deletedAt gorm.DeletedAt
	deletedAt.Time = time.Now()
	deletedAt.Valid = true
	if res := dao.GormDB.Model(&model.Session{}).Where("send_id = ? AND receive_id = ?", req.UserId, req.GroupId).Update("deleted_at", deletedAt); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除联系人
	if res := dao.GormDB.Model(&model.UserContact{}).Where("user_id = ? AND contact_id = ?", req.UserId, req.GroupId).Updates(map[string]interface{}{
		"deleted_at": deletedAt,
		"status":     enum.QUIT_GROUP, // 退群
	}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除申请记录，后面还可以加
	if res := dao.GormDB.Model(&model.ContactApply{}).Where("contact_id = ? AND user_id = ?", req.GroupId, req.UserId).Update("deleted_at", deletedAt); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	//if err := myredis.DelKeysWithPattern("group_info_" + groupId); err != nil {
	//	zlog.Error(err.Error())
	//}
	//if err := myredis.DelKeysWithPattern("groupmember_list_" + groupId); err != nil {
	//	zlog.Error(err.Error())
	//}
	if err := myredis.DelKeysWithPattern("group_session_list_" + req.UserId); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPattern("my_joined_group_list_ " + req.UserId); err != nil {
		zlog.Error(err.Error())
	}
	//if err := myredis.DelKeysWithPattern("session_" + userId + "_" + groupId); err != nil {
	//	zlog.Error(err.Error())
	//}
	return "退群成功", 0
}

// DismissGroup 解散群聊
func (gis *GroupInfoService) DismissGroup(req *request.DismissGroupRequest) (string, int) {
	var deletedAt gorm.DeletedAt
	deletedAt.Time = time.Now()
	deletedAt.Valid = true
	if res := dao.GormDB.Model(&model.GroupInfo{}).Where("uuid = ?", req.GroupId).Updates(
		map[string]interface{}{
			"deleted_at": deletedAt,
			"updated_at": deletedAt.Time,
		}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	var sessionList []model.Session
	if res := dao.GormDB.Model(&model.Session{}).Where("receive_id = ?", req.GroupId).Find(&sessionList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, session := range sessionList {
		if res := dao.GormDB.Model(&session).Updates(
			map[string]interface{}{
				"deleted_at": deletedAt,
			}); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}

	var userContactList []model.UserContact
	if res := dao.GormDB.Model(&model.UserContact{}).Where("contact_id = ?", req.GroupId).Find(&userContactList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	for _, userContact := range userContactList {
		if res := dao.GormDB.Model(&userContact).Update("deleted_at", deletedAt); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}

	var contactApplys []model.ContactApply
	if res := dao.GormDB.Model(&contactApplys).Where("contact_id = ?", req.GroupId).Find(&contactApplys); res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			zlog.Info(res.Error.Error())
			return "无响应的申请记录需要删除", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, contactApply := range contactApplys {
		if res := dao.GormDB.Model(&contactApply).Update("deleted_at", deletedAt); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}
	//if err := myredis.DelKeysWithPattern("group_info_" + groupId); err != nil {
	//	zlog.Error(err.Error())
	//}
	//if err := myredis.DelKeysWithPattern("groupmember_list_" + groupId); err != nil {
	//	zlog.Error(err.Error())
	//}
	if err := myredis.DelKeysWithPattern("contact_mygroup_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPattern("group_session_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPrefix("my_joined_group_list"); err != nil {
		zlog.Error(err.Error())
	}
	return "解散群聊成功", 0
}

// GetGroupInfo 获取群聊详情
func (gis *GroupInfoService) GetGroupInfo(req *request.GetGroupInfoRequest) (string, *respond.GetGroupInfoRespond, int) {
	rspString, err := myredis.GetKeyNilIsErr("group_info_" + req.GroupId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var group model.GroupInfo
			if res := dao.GormDB.First(&group, "uuid = ?", req.GroupId); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, nil, -1
			}
			rsp := &respond.GetGroupInfoRespond{
				Uuid:      group.Uuid,
				Name:      group.Name,
				Notice:    group.Notice,
				Avatar:    group.Avatar,
				MemberCnt: group.MemberCnt,
				OwnerId:   group.OwnerId,
				AddMode:   group.AddMode,
				Status:    group.Status,
			}
			if group.DeletedAt.Valid {
				rsp.IsDeleted = true
			} else {
				rsp.IsDeleted = false
			}
			//rspString, err := json.Marshal(rsp)
			//if err != nil {
			//	zlog.Error(err.Error())
			//}
			//if err := myredis.SetKeyEx("group_info_"+groupId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
			//	zlog.Error(err.Error())
			//}
			return "获取成功", rsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	var rsp *respond.GetGroupInfoRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取成功", rsp, 0
}

// GetGroupInfoList 获取群聊列表 - 管理员
// 管理员少，而且如果用户更改了，那么管理员会一直频繁删除redis，更新redis，比较麻烦，所以管理员暂时不使用redis缓存
func (gis *GroupInfoService) GetGroupInfoList() (string, []respond.GetGroupListRespond, int) {
	var groupList []model.GroupInfo
	if res := dao.GormDB.Unscoped().Find(&groupList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var rsp []respond.GetGroupListRespond
	for _, group := range groupList {
		rp := respond.GetGroupListRespond{
			Uuid:    group.Uuid,
			Name:    group.Name,
			OwnerId: group.OwnerId,
			Status:  group.Status,
		}
		if group.DeletedAt.Valid {
			rp.IsDeleted = true
		} else {
			rp.IsDeleted = false
		}
		rsp = append(rsp, rp)
	}
	return "获取成功", rsp, 0
}

// DeleteGroups 删除列表中群聊 - 管理员
func (gis *GroupInfoService) DeleteGroups(req *request.DeleteGroupsRequest) (string, int) {
	for _, uuid := range req.UuidList {
		var deletedAt gorm.DeletedAt
		deletedAt.Time = time.Now()
		deletedAt.Valid = true
		if res := dao.GormDB.Model(&model.GroupInfo{}).Where("uuid = ?", uuid).Update("deleted_at", deletedAt); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		// 删除会话
		var sessionList []model.Session
		if res := dao.GormDB.Model(&model.Session{}).Where("receive_id = ?", uuid).Find(&sessionList); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		for _, session := range sessionList {
			if res := dao.GormDB.Model(&session).Update("deleted_at", deletedAt); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
		// 删除联系人
		var userContactList []model.UserContact
		if res := dao.GormDB.Model(&model.UserContact{}).Where("contact_id = ?", uuid).Find(&userContactList); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}

		for _, userContact := range userContactList {
			if res := dao.GormDB.Model(&userContact).Update("deleted_at", deletedAt); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}

		var contactApplys []model.ContactApply
		if res := dao.GormDB.Model(&contactApplys).Where("contact_id = ?", uuid).Find(&contactApplys); res.Error != nil {
			if res.Error != gorm.ErrRecordNotFound {
				zlog.Info(res.Error.Error())
				return "无响应的申请记录需要删除", 0
			}
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
		for _, contactApply := range contactApplys {
			if res := dao.GormDB.Model(&contactApply).Update("deleted_at", deletedAt); res.Error != nil {
				zlog.Error(res.Error.Error())
				return constants.SYSTEM_ERROR, -1
			}
		}
	}
	//for _, uuid := range uuidList {
	//	if err := myredis.DelKeysWithPattern("group_info_" + uuid); err != nil {
	//		zlog.Error(err.Error())
	//	}
	//	if err := myredis.DelKeysWithPattern("groupmember_list_" + uuid); err != nil {
	//		zlog.Error(err.Error())
	//	}
	//}
	if err := myredis.DelKeysWithPrefix("contact_mygroup_list"); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPrefix("group_session_list"); err != nil {
		zlog.Error(err.Error())
	}
	if err := myredis.DelKeysWithPrefix("group_session_list"); err != nil {
		zlog.Error(err.Error())
	}
	return "解散/删除群聊成功", 0
}
