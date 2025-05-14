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

type SessionService struct {
	Ctx *gin.Context
}

// CreateSession 创建会话
func (ss *SessionService) CreateSession(req *request.CreateSessionRequest) (string, string, int) {
	var user model.UserInfo
	// 获取用户信息
	if res := dao.GormDB.Where("uuid = ?", req.SendId).First(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, "", -1
	}
	// 创建session会话
	var session model.Session
	session.Uuid = fmt.Sprintf("S%s", random.GetNowAndLenRandomString(11))
	session.SendId = req.SendId
	session.ReceiveId = req.ReceiveId
	session.CreatedAt = time.Now()
	// req.ReceiveId[0] == 'U' 代表是用户对话
	if req.ReceiveId[0] == 'U' {
		var receiveUser model.UserInfo
		// 获取接收用户信息
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&receiveUser); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, "", -1
		}
		// 判断用户是否被禁用
		if receiveUser.Status == enum.DISABLE {
			zlog.Error("该用户被禁用了")
			return "该用户被禁用了", "", -2
		} else {
			session.ReceiveName = receiveUser.Nickname
			session.Avatar = receiveUser.Avatar
		}
	} else { // req.ReceiveId[0] == 'G' 代表是群聊对话
		var receiveGroup model.GroupInfo
		// 获取接收群聊信息
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&receiveGroup); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, "", -1
		}
		// 判断群聊是否被禁用
		if receiveGroup.Status == enum.DISABLE {
			zlog.Error("该群聊被禁用了")
			return "该群聊被禁用了", "", -2
		} else {
			session.ReceiveName = receiveGroup.Name
			session.Avatar = receiveGroup.Avatar
		}
	}
	// 创建会话
	if res := dao.GormDB.Create(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, "", -1
	}
	// 删除群组redis
	if err := myredis.DelKeysWithPattern("group_session_list_" + req.SendId); err != nil {
		zlog.Error(err.Error())
	}
	// 删除会话redis
	if err := myredis.DelKeysWithPattern("session_list_" + req.SendId); err != nil {
		zlog.Error(err.Error())
	}
	return "会话创建成功", session.Uuid, 0
}

// CheckOpenSessionAllowed 检查是否允许发起会话
func (ss *SessionService) CheckOpenSessionAllowed(req *request.CreateSessionRequest) (string, bool, int) {
	var contact model.UserContact
	// 判断是否允许发起会话
	if res := dao.GormDB.Where("user_id = ? and contact_id = ?", req.SendId, req.ReceiveId).First(&contact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, false, -1
	}
	// 判断是否可以发起会话
	if contact.Status == enum.BE_BLACK {
		return "已被对方拉黑，无法发起会话", false, -2
	} else if contact.Status == enum.BLACK {
		return "已拉黑对方，先解除拉黑状态才能发起会话", false, -2
	}
	// req.ReceiveId[0] == 'U' 代表是用户对话
	if req.ReceiveId[0] == 'U' {
		var user model.UserInfo
		// 获取接收用户信息
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&user); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, false, -1
		}
		if user.Status == enum.DISABLE {
			zlog.Info("对方已被禁用，无法发起会话")
			return "对方已被禁用，无法发起会话", false, -2
		}
	} else { // req.ReceiveId[0] == 'G' 代表是群聊对话
		var group model.GroupInfo
		// 获取接收群聊信息
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&group); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, false, -1
		}
		if group.Status == enum.DISABLE {
			zlog.Info("对方已被禁用，无法发起会话")
			return "对方已被禁用，无法发起会话", false, -2
		}
	}
	return "可以发起会话", true, 0
}

// OpenSession 打开会话
func (ss *SessionService) OpenSession(req *request.OpenSessionRequest) (string, string, int) {
	// 从redis获取会话
	rspString, err := myredis.GetKeyWithPrefixNilIsErr("session_" + req.SendId + "_" + req.ReceiveId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 从数据库中获取会话
			var session model.Session
			if res := dao.GormDB.Where("send_id = ? and receive_id = ?", req.SendId, req.ReceiveId).First(&session); res.Error != nil {
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					// 会话没有找到，将新建会话
					zlog.Info("会话没有找到，将新建会话")
					createReq := &request.CreateSessionRequest{
						SendId:    req.SendId,
						ReceiveId: req.ReceiveId,
					}
					return ss.CreateSession(createReq)
				}
			}
			// 将会话写入redis
			rspString, err := json.Marshal(session)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("session_"+req.SendId+"_"+req.ReceiveId+"_"+session.Uuid, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "会话创建成功", session.Uuid, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, "", -1
		}
	}
	// 将redis中的会话转为结构体
	var session model.Session
	if err := json.Unmarshal([]byte(rspString), &session); err != nil {
		zlog.Error(err.Error())
	}
	return "会话创建成功", session.Uuid, 0
}

// GetUserSessionList 获取用户会话列表
func (ss *SessionService) GetUserSessionList(req *request.OwnlistRequest) (string, []respond.UserSessionListRespond, int) {
	// 从redis获取用户会话列表
	rspString, err := myredis.GetKeyNilIsErr("session_list_" + req.OwnerId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var sessionList []model.Session
			// 获取会话列表
			if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", req.OwnerId).Find(&sessionList); res.Error != nil {
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					zlog.Info("未创建用户会话")
					return "未创建用户会话", nil, 0
				} else {
					zlog.Error(res.Error.Error())
					return constants.SYSTEM_ERROR, nil, -1
				}
			}
			// 将会话列表写入redis
			var sessionListRsp []respond.UserSessionListRespond
			for i := 0; i < len(sessionList); i++ {
				if sessionList[i].ReceiveId[0] == 'U' {
					sessionListRsp = append(sessionListRsp, respond.UserSessionListRespond{
						SessionId: sessionList[i].Uuid,
						Avatar:    sessionList[i].Avatar,
						UserId:    sessionList[i].ReceiveId,
						Username:  sessionList[i].ReceiveName,
					})
				}
			}
			// 将会话列表转为json
			rspString, err := json.Marshal(sessionListRsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("session_list_"+req.OwnerId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取成功", sessionListRsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	// 将redis中的会话列表转为结构体
	var rsp []respond.UserSessionListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取成功", rsp, 0
}

// GetGroupSessionList 获取群聊会话列表
func (ss *SessionService) GetGroupSessionList(req *request.OwnlistRequest) (string, []respond.GroupSessionListRespond, int) {
	// 从redis获取群聊会话列表
	rspString, err := myredis.GetKeyNilIsErr("group_session_list_" + req.OwnerId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var sessionList []model.Session
			// 获取会话列表
			if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", req.OwnerId).Find(&sessionList); res.Error != nil {
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					zlog.Info("未创建群聊会话")
					return "未创建群聊会话", nil, 0
				} else {
					zlog.Error(res.Error.Error())
					return constants.SYSTEM_ERROR, nil, -1
				}
			}
			// 将会话列表写入redis
			var sessionListRsp []respond.GroupSessionListRespond
			for i := 0; i < len(sessionList); i++ {
				if sessionList[i].ReceiveId[0] == 'G' {
					sessionListRsp = append(sessionListRsp, respond.GroupSessionListRespond{
						SessionId: sessionList[i].Uuid,
						Avatar:    sessionList[i].Avatar,
						GroupId:   sessionList[i].ReceiveId,
						GroupName: sessionList[i].ReceiveName,
					})
				}
			}
			// 将会话列表转为json
			rspString, err := json.Marshal(sessionListRsp)
			if err != nil {
				zlog.Error(err.Error())
			}
			if err := myredis.SetKeyEx("group_session_list_"+req.OwnerId, string(rspString), time.Minute*constants.REDIS_TIMEOUT); err != nil {
				zlog.Error(err.Error())
			}
			return "获取成功", sessionListRsp, 0
		} else {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, nil, -1
		}
	}
	// 将redis中的会话列表转为结构体
	var rsp []respond.GroupSessionListRespond
	if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
		zlog.Error(err.Error())
	}
	return "获取成功", rsp, 0
}

// DeleteSession 删除会话
func (ss *SessionService) DeleteSession(req *request.DeleteSessionRequest) (string, int) {
	var session model.Session
	// 从数据库中获取会话
	if res := dao.GormDB.Where("uuid = ?", req.SessionId).Find(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	session.DeletedAt.Valid = true
	session.DeletedAt.Time = time.Now()
	if res := dao.GormDB.Save(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除会话
	if err := myredis.DelKeysWithSuffix(req.SessionId); err != nil {
		zlog.Error(err.Error())
	}
	// 删除群聊会话列表
	if err := myredis.DelKeysWithPattern("group_session_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	// 删除用户会话列表
	if err := myredis.DelKeysWithPattern("session_list_" + req.OwnerId); err != nil {
		zlog.Error(err.Error())
	}
	return "删除成功", 0
}
