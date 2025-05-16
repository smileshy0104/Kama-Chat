package chat

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
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
	"time"
)

// Server 定义聊天服务器的结构体
// 用于管理客户端连接、消息转发以及客户端登录/登出等操作
type Server struct {
	// Clients 存储所有在线客户端，以客户端 UUID 为键，*Client 对象为值
	Clients map[string]*Client
	// mutex 用于保护 Clients 映射的并发访问，确保线程安全
	mutex *sync.Mutex
	// Transmit 消息转发通道，用于将接收到的消息广播给所有在线客户端
	Transmit chan []byte // 转发通道
	// Login 登录通道，接收新上线的客户端对象，用于添加到在线列表
	Login chan *Client // 登录通道
	// Logout 登出通道，接收下线的客户端对象，用于从在线列表中移除
	Logout chan *Client // 退出登录通道
}

// ChatServer 是 Server 的全局实例，表示当前运行的聊天服务器
var ChatServer *Server

// init函数用于初始化ChatServer实例。
// 当ChatServer实例不存在时，通过创建一个新的Server实例并将其赋值给ChatServer。
// 这个初始化过程包括：
// - 创建一个空的Clients字典，用于后续存储客户端信息。
// - 初始化一个互斥锁mutex，用于在并发环境下保护Clients字典的安全访问。
// - 创建Transmit、Login和Logout通道，分别用于消息传输、客户端登录和注销的管理。
func init() {
	// 如果 ChatServer 尚未初始化，则创建一个新实例
	if ChatServer == nil {
		// 创建一个新实例
		ChatServer = &Server{
			Clients:  make(map[string]*Client),                   // 创建一个空的Clients字典
			mutex:    &sync.Mutex{},                              // 创建一个互斥锁
			Transmit: make(chan []byte, constants.CHANNEL_SIZE),  // 创建一个Transmit通道
			Login:    make(chan *Client, constants.CHANNEL_SIZE), // 创建一个Login通道
			Logout:   make(chan *Client, constants.CHANNEL_SIZE), // 创建一个Logout通道
		}
	}
}

// 将https://127.0.0.1:8000/static/xxx 转为 /static/xxx
func normalizePath(path string) string {
	// 查找 "/static/" 的位置
	if path == "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" {
		return path
	}
	staticIndex := strings.Index(path, "/static/")
	if staticIndex < 0 {
		log.Println(path)
		zlog.Error("路径不合法")
	}
	// 返回从 "/static/" 开始的部分
	return path[staticIndex:]
}

// Start 启动函数，Server端用主进程起，Client端可以用协程起
func (s *Server) Start() {
	defer func() {
		close(s.Transmit)
		close(s.Logout)
		close(s.Login)
	}()
	for {
		select {
		case client := <-s.Login:
			{
				s.mutex.Lock()
				s.Clients[client.Uuid] = client
				s.mutex.Unlock()
				zlog.Debug(fmt.Sprintf("欢迎来到kama聊天服务器，亲爱的用户%s\n", client.Uuid))
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte("欢迎来到kama聊天服务器"))
				if err != nil {
					zlog.Error(err.Error())
				}
			}

		case client := <-s.Logout:
			{
				s.mutex.Lock()
				delete(s.Clients, client.Uuid)
				s.mutex.Unlock()
				zlog.Info(fmt.Sprintf("用户%s退出登录\n", client.Uuid))
				if err := client.Conn.WriteMessage(websocket.TextMessage, []byte("已退出登录")); err != nil {
					zlog.Error(err.Error())
				}
			}

		case data := <-s.Transmit:
			{
				var chatMessageReq request.ChatMessageRequest
				if err := json.Unmarshal(data, &chatMessageReq); err != nil {
					zlog.Error(err.Error())
				}
				// log.Println("原消息为：", data, "反序列化后为：", chatMessageReq)
				if chatMessageReq.Type == enum.Text {
					// 存message
					message := model.Message{
						Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
						SessionId:  chatMessageReq.SessionId,
						Type:       chatMessageReq.Type,
						Content:    chatMessageReq.Content,
						Url:        "",
						SendId:     chatMessageReq.SendId,
						SendName:   chatMessageReq.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  chatMessageReq.ReceiveId,
						FileSize:   "0B",
						FileType:   "",
						FileName:   "",
						Status:     enum.Unsent,
						CreatedAt:  time.Now(),
						AVdata:     "",
					}
					// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
					message.SendAvatar = normalizePath(message.SendAvatar)
					if res := dao.GormDB.Create(&message); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					if message.ReceiveId[0] == 'U' { // 发送给User
						// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
						// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
						// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
						messageRsp := respond.GetMessageListRespond{
							SendId:     message.SendId,
							SendName:   message.SendName,
							SendAvatar: chatMessageReq.SendAvatar,
							ReceiveId:  message.ReceiveId,
							Type:       message.Type,
							Content:    message.Content,
							Url:        message.Url,
							FileSize:   message.FileSize,
							FileName:   message.FileName,
							FileType:   message.FileType,
							CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
						}
						jsonMessage, err := json.Marshal(messageRsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
						var messageBack = &MessageBack{
							Message: jsonMessage,
							Uuid:    message.Uuid,
						}
						s.mutex.Lock()
						if receiveClient, ok := s.Clients[message.ReceiveId]; ok {
							//messageBack.Message = jsonMessage
							//messageBack.Uuid = message.Uuid
							receiveClient.SendBack <- messageBack // 向client.Send发送
						}
						// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
						// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
						// 所以这里后端进行回显，前端不回显
						sendClient := s.Clients[message.SendId]
						sendClient.SendBack <- messageBack
						s.mutex.Unlock()

						// redis
						var rspString string
						rspString, err = myredis.GetKeyNilIsErr("message_list_" + message.SendId + "_" + message.ReceiveId)
						if err == nil {
							var rsp []respond.GetMessageListRespond
							if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
								zlog.Error(err.Error())
							}
							rsp = append(rsp, messageRsp)
							rspByte, err := json.Marshal(rsp)
							if err != nil {
								zlog.Error(err.Error())
							}
							if err := myredis.SetKeyEx("message_list_"+message.SendId+"_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
								zlog.Error(err.Error())
							}
						} else {
							if !errors.Is(err, redis.Nil) {
								zlog.Error(err.Error())
							}
						}

					} else if message.ReceiveId[0] == 'G' { // 发送给Group
						messageRsp := respond.GetGroupMessageListRespond{
							SendId:     message.SendId,
							SendName:   message.SendName,
							SendAvatar: chatMessageReq.SendAvatar,
							ReceiveId:  message.ReceiveId,
							Type:       message.Type,
							Content:    message.Content,
							Url:        message.Url,
							FileSize:   message.FileSize,
							FileName:   message.FileName,
							FileType:   message.FileType,
							CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
						}
						jsonMessage, err := json.Marshal(messageRsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
						var messageBack = &MessageBack{
							Message: jsonMessage,
							Uuid:    message.Uuid,
						}
						var group model.GroupInfo
						if res := dao.GormDB.Where("uuid = ?", message.ReceiveId).First(&group); res.Error != nil {
							zlog.Error(res.Error.Error())
						}
						var members []string
						if err := json.Unmarshal(group.Members, &members); err != nil {
							zlog.Error(err.Error())
						}
						s.mutex.Lock()
						for _, member := range members {
							if member != message.SendId {
								if receiveClient, ok := s.Clients[member]; ok {
									receiveClient.SendBack <- messageBack
								}
							} else {
								sendClient := s.Clients[message.SendId]
								sendClient.SendBack <- messageBack
							}
						}
						s.mutex.Unlock()

						// redis
						var rspString string
						rspString, err = myredis.GetKeyNilIsErr("group_messagelist_" + message.ReceiveId)
						if err == nil {
							var rsp []respond.GetGroupMessageListRespond
							if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
								zlog.Error(err.Error())
							}
							rsp = append(rsp, messageRsp)
							rspByte, err := json.Marshal(rsp)
							if err != nil {
								zlog.Error(err.Error())
							}
							if err := myredis.SetKeyEx("group_messagelist_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
								zlog.Error(err.Error())
							}
						} else {
							if !errors.Is(err, redis.Nil) {
								zlog.Error(err.Error())
							}
						}
					}
				} else if chatMessageReq.Type == enum.File {
					// 存message
					message := model.Message{
						Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
						SessionId:  chatMessageReq.SessionId,
						Type:       chatMessageReq.Type,
						Content:    "",
						Url:        chatMessageReq.Url,
						SendId:     chatMessageReq.SendId,
						SendName:   chatMessageReq.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  chatMessageReq.ReceiveId,
						FileSize:   chatMessageReq.FileSize,
						FileType:   chatMessageReq.FileType,
						FileName:   chatMessageReq.FileName,
						Status:     enum.Unsent,
						CreatedAt:  time.Now(),
						AVdata:     "",
					}
					// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
					message.SendAvatar = normalizePath(message.SendAvatar)
					if res := dao.GormDB.Create(&message); res.Error != nil {
						zlog.Error(res.Error.Error())
					}
					if message.ReceiveId[0] == 'U' { // 发送给User
						// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
						// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
						// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
						messageRsp := respond.GetMessageListRespond{
							SendId:     message.SendId,
							SendName:   message.SendName,
							SendAvatar: chatMessageReq.SendAvatar,
							ReceiveId:  message.ReceiveId,
							Type:       message.Type,
							Content:    message.Content,
							Url:        message.Url,
							FileSize:   message.FileSize,
							FileName:   message.FileName,
							FileType:   message.FileType,
							CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
						}
						jsonMessage, err := json.Marshal(messageRsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
						var messageBack = &MessageBack{
							Message: jsonMessage,
							Uuid:    message.Uuid,
						}
						s.mutex.Lock()
						if receiveClient, ok := s.Clients[message.ReceiveId]; ok {
							//messageBack.Message = jsonMessage
							//messageBack.Uuid = message.Uuid
							receiveClient.SendBack <- messageBack // 向client.Send发送
						}
						// 因为send_id肯定在线，所以这里在后端进行在线回显message，其实优化的话前端可以直接回显
						// 问题在于前后端的req和rsp结构不同，前端存储message的messageList不能存req，只能存rsp
						// 所以这里后端进行回显，前端不回显
						sendClient := s.Clients[message.SendId]
						sendClient.SendBack <- messageBack
						s.mutex.Unlock()

						// redis
						var rspString string
						rspString, err = myredis.GetKeyNilIsErr("message_list_" + message.SendId + "_" + message.ReceiveId)
						if err == nil {
							var rsp []respond.GetMessageListRespond
							if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
								zlog.Error(err.Error())
							}
							rsp = append(rsp, messageRsp)
							rspByte, err := json.Marshal(rsp)
							if err != nil {
								zlog.Error(err.Error())
							}
							if err := myredis.SetKeyEx("message_list_"+message.SendId+"_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
								zlog.Error(err.Error())
							}
						} else {
							if !errors.Is(err, redis.Nil) {
								zlog.Error(err.Error())
							}
						}
					} else {
						messageRsp := respond.GetGroupMessageListRespond{
							SendId:     message.SendId,
							SendName:   message.SendName,
							SendAvatar: chatMessageReq.SendAvatar,
							ReceiveId:  message.ReceiveId,
							Type:       message.Type,
							Content:    message.Content,
							Url:        message.Url,
							FileSize:   message.FileSize,
							FileName:   message.FileName,
							FileType:   message.FileType,
							CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
						}
						jsonMessage, err := json.Marshal(messageRsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
						var messageBack = &MessageBack{
							Message: jsonMessage,
							Uuid:    message.Uuid,
						}
						var group model.GroupInfo
						if res := dao.GormDB.Where("uuid = ?", message.ReceiveId).First(&group); res.Error != nil {
							zlog.Error(res.Error.Error())
						}
						var members []string
						if err := json.Unmarshal(group.Members, &members); err != nil {
							zlog.Error(err.Error())
						}
						s.mutex.Lock()
						for _, member := range members {
							if member != message.SendId {
								if receiveClient, ok := s.Clients[member]; ok {
									receiveClient.SendBack <- messageBack
								}
							} else {
								sendClient := s.Clients[message.SendId]
								sendClient.SendBack <- messageBack
							}
						}
						s.mutex.Unlock()

						// redis
						var rspString string
						rspString, err = myredis.GetKeyNilIsErr("group_messagelist_" + message.ReceiveId)
						if err == nil {
							var rsp []respond.GetGroupMessageListRespond
							if err := json.Unmarshal([]byte(rspString), &rsp); err != nil {
								zlog.Error(err.Error())
							}
							rsp = append(rsp, messageRsp)
							rspByte, err := json.Marshal(rsp)
							if err != nil {
								zlog.Error(err.Error())
							}
							if err := myredis.SetKeyEx("group_messagelist_"+message.ReceiveId, string(rspByte), time.Minute*constants.REDIS_TIMEOUT); err != nil {
								zlog.Error(err.Error())
							}
						} else {
							if !errors.Is(err, redis.Nil) {
								zlog.Error(err.Error())
							}
						}
					}
				} else if chatMessageReq.Type == enum.AudioOrVideo {
					var avData request.AVData
					if err := json.Unmarshal([]byte(chatMessageReq.AVdata), &avData); err != nil {
						zlog.Error(err.Error())
					}
					//log.Println(avData)
					message := model.Message{
						Uuid:       fmt.Sprintf("M%s", random.GetNowAndLenRandomString(11)),
						SessionId:  chatMessageReq.SessionId,
						Type:       chatMessageReq.Type,
						Content:    "",
						Url:        "",
						SendId:     chatMessageReq.SendId,
						SendName:   chatMessageReq.SendName,
						SendAvatar: chatMessageReq.SendAvatar,
						ReceiveId:  chatMessageReq.ReceiveId,
						FileSize:   "",
						FileType:   "",
						FileName:   "",
						Status:     enum.Unsent,
						CreatedAt:  time.Now(),
						AVdata:     chatMessageReq.AVdata,
					}
					if avData.MessageId == "PROXY" && (avData.Type == "start_call" || avData.Type == "receive_call" || avData.Type == "reject_call") {
						// 存message
						// 对SendAvatar去除前面/static之前的所有内容，防止ip前缀引入
						message.SendAvatar = normalizePath(message.SendAvatar)
						if res := dao.GormDB.Create(&message); res.Error != nil {
							zlog.Error(res.Error.Error())
						}
					}

					if chatMessageReq.ReceiveId[0] == 'U' { // 发送给User
						// 如果能找到ReceiveId，说明在线，可以发送，否则存表后跳过
						// 因为在线的时候是通过websocket更新消息记录的，离线后通过存表，登录时只调用一次数据库操作
						// 切换chat对象后，前端的messageList也会改变，获取messageList从第二次就是从redis中获取
						messageRsp := respond.AVMessageRespond{
							SendId:     message.SendId,
							SendName:   message.SendName,
							SendAvatar: message.SendAvatar,
							ReceiveId:  message.ReceiveId,
							Type:       message.Type,
							Content:    message.Content,
							Url:        message.Url,
							FileSize:   message.FileSize,
							FileName:   message.FileName,
							FileType:   message.FileType,
							CreatedAt:  message.CreatedAt.Format("2006-01-02 15:04:05"),
							AVdata:     message.AVdata,
						}
						jsonMessage, err := json.Marshal(messageRsp)
						if err != nil {
							zlog.Error(err.Error())
						}
						// log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)
						log.Println("返回的消息为：", messageRsp)
						var messageBack = &MessageBack{
							Message: jsonMessage,
							Uuid:    message.Uuid,
						}
						s.mutex.Lock()
						if receiveClient, ok := s.Clients[message.ReceiveId]; ok {
							//messageBack.Message = jsonMessage
							//messageBack.Uuid = message.Uuid
							receiveClient.SendBack <- messageBack // 向client.Send发送
						}
						// 通话这不能回显，发回去的话就会出现两个start_call。
						//sendClient := s.Clients[message.SendId]
						//sendClient.SendBack <- messageBack
						s.mutex.Unlock()
					}
				}

			}
		}
	}
}

func (s *Server) Close() {
	close(s.Login)
	close(s.Logout)
	close(s.Transmit)
}

func (s *Server) SendClientToLogin(client *Client) {
	s.mutex.Lock()
	s.Login <- client
	s.mutex.Unlock()
}

func (s *Server) SendClientToLogout(client *Client) {
	s.mutex.Lock()
	s.Logout <- client
	s.mutex.Unlock()
}

func (s *Server) SendMessageToTransmit(message []byte) {
	s.mutex.Lock()
	s.Transmit <- message
	s.mutex.Unlock()
}

func (s *Server) RemoveClient(uuid string) {
	s.mutex.Lock()
	delete(s.Clients, uuid)
	s.mutex.Unlock()
}
