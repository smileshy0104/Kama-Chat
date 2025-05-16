package chat

import (
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/lib/kafka"
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
	"os"
	"sync"
	"time"
)

// KafkaServer 定义了基于 Kafka 的服务器结构体，用于管理客户端连接以及登录/登出事件。
type KafkaServer struct {
	// Clients 存储所有当前连接的客户端，以客户端 UUID 为键，*Client 对象为值。
	Clients map[string]*Client
	// mutex 用于确保对 Clients 映射的并发访问是线程安全的。
	mutex *sync.Mutex
	// Login 登录通道，用于通知有新的客户端登录事件。
	Login chan *Client
	// Logout 退出登录通道，用于通知有客户端退出登录。
	Logout chan *Client
}

// KafkaChatServer 是 KafkaServer 的全局实例，用于管理所有客户端连接。
var KafkaChatServer *KafkaServer

// kafkaQuit 是一个用于通知 Kafka 服务器退出的通道，用于处理终止信号。
var kafkaQuit = make(chan os.Signal, 1)

// init函数用于初始化KafkaChatServer实例。
// 该函数确保在程序启动时，如果KafkaChatServer尚未被初始化，则将其初始化。
// 这是因为KafkaChatServer可能在整个应用程序中被多个部分使用，而初始化应该只发生一次。
func InitKafka() {
	// 检查KafkaChatServer是否已经初始化。
	if KafkaChatServer == nil {
		// 如果没有初始化，则创建一个新的KafkaServer实例。
		// 这包括创建一个空的Clients映射，用于跟踪当前在线的客户端，
		// 以及初始化用于同步访问的互斥锁和用于通知的通道。
		KafkaChatServer = &KafkaServer{
			Clients: make(map[string]*Client),
			mutex:   &sync.Mutex{},
			Login:   make(chan *Client),
			Logout:  make(chan *Client),
		}
	}
	// 初始化一个信号通知，当接收到中断或终止信号时，kafkaQuit通道将被用来优雅地关闭服务。
	// 注意：此行代码已被注释掉，可能是因为在当前上下文中不需要信号处理，或者信号处理已经在其他地方实现。
	// signal.Notify(kafkaQuit, syscall.SIGINT, syscall.SIGTERM)
}

// Start 方法用于启动 Kafka 服务器。
// 它负责初始化服务器、监听传入的连接，并处理来自 Kafka 的消息以及客户端的登录/登出事件。

func (k *KafkaServer) Start() {
	// 在函数退出时关闭 Login 和 Logout 通道，确保资源释放。
	defer func() {
		if r := recover(); r != nil {
			zlog.Error(fmt.Sprintf("kafka server panic: %v", r))
		}
		close(k.Login)
		close(k.Logout)
	}()

	// 启动一个 goroutine 来读取 Kafka 中的消息。
	go func() {
		// 捕获可能发生的 panic 并记录日志。
		defer func() {
			if r := recover(); r != nil {
				zlog.Error(fmt.Sprintf("kafka server panic: %v", r))
			}
		}()

		// 进入循环不断从 Kafka 读取消息。
		for {
			kafkaMessage, err := kafka.KafkaService.ChatReader.ReadMessage(ctx)
			if err != nil {
				zlog.Error(err.Error())
			}

			// 记录 Kafka 消息的详细信息。
			log.Printf("topic=%s, partition=%d, offset=%d, key=%s, value=%s",
				kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset,
				kafkaMessage.Key, kafkaMessage.Value)
			zlog.Info(fmt.Sprintf("topic=%s, partition=%d, offset=%d, key=%s, value=%s",
				kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset,
				kafkaMessage.Key, kafkaMessage.Value))

			// 解析接收到的消息内容为 ChatMessageRequest 结构体。
			data := kafkaMessage.Value
			var chatMessageReq request.ChatMessageRequest
			if err := json.Unmarshal(data, &chatMessageReq); err != nil {
				zlog.Error(err.Error())
			}
			log.Println("原消息为：", data, "反序列化后为：", chatMessageReq)

			// 根据消息类型进行处理。
			if chatMessageReq.Type == enum.Text { // 处理文本消息
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

				// 对 SendAvatar 去除前面 "/static" 之前的内容，防止 IP 前缀引入。
				message.SendAvatar = normalizePath(message.SendAvatar)

				// 将消息保存到数据库。
				if res := dao.GormDB.Create(&message); res.Error != nil {
					zlog.Error(res.Error.Error())
				}

				// 判断接收者是用户还是群组。
				if message.ReceiveId[0] == 'U' { // 发送给用户
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

					// 序列化消息以便发送。
					jsonMessage, err := json.Marshal(messageRsp)
					if err != nil {
						zlog.Error(err.Error())
					}
					log.Println("返回的消息为：", messageRsp, "序列化后为：", jsonMessage)

					messageBack := &MessageBack{
						Message: jsonMessage,
						Uuid:    message.Uuid,
					}

					k.mutex.Lock()
					if receiveClient, ok := k.Clients[message.ReceiveId]; ok {
						receiveClient.SendBack <- messageBack
					}

					// 因为发送方一定在线，所以直接回显消息给发送方。
					sendClient := k.Clients[message.SendId]
					sendClient.SendBack <- messageBack
					k.mutex.Unlock()

					// 更新 Redis 缓存。
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
					} else if !errors.Is(err, redis.Nil) {
						zlog.Error(err.Error())
					}

				} else if message.ReceiveId[0] == 'G' { // 发送给群组
					// 群组消息处理逻辑与用户类似，但需要遍历所有群组成员并发送消息。
					// ...
				}
			} else if chatMessageReq.Type == enum.File { // 处理文件消息
				// 文件消息的处理逻辑类似于文本消息，包含文件相关信息。
				// ...
			} else if chatMessageReq.Type == enum.AudioOrVideo { // 处理音视频消息
				// 音视频消息（如通话请求）的处理逻辑。
				// ...
			}
		}
	}()

	// 主循环监听客户端的登录和登出事件。
	for {
		select {
		case client := <-k.Login:
			{
				// 添加新登录的客户端。
				k.mutex.Lock()
				k.Clients[client.Uuid] = client
				k.mutex.Unlock()
				zlog.Debug(fmt.Sprintf("欢迎来到 Kama 聊天服务器，亲爱的用户 %s\n", client.Uuid))
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte("欢迎来到 Kama 聊天服务器"))
				if err != nil {
					zlog.Error(err.Error())
				}
			}

		case client := <-k.Logout:
			{
				// 移除已登出的客户端。
				k.mutex.Lock()
				delete(k.Clients, client.Uuid)
				k.mutex.Unlock()
				zlog.Info(fmt.Sprintf("用户 %s 退出登录\n", client.Uuid))
				if err := client.Conn.WriteMessage(websocket.TextMessage, []byte("已退出登录")); err != nil {
					zlog.Error(err.Error())
				}
			}
		}
	}
}

// Close 关闭Kafka服务器的登录和登出通道。
// 该方法旨在清理Kafka服务器的资源，确保服务安全停止。
func (k *KafkaServer) Close() {
	// 关闭登录和登出通道，确保服务安全停止。
	close(k.Login)
	close(k.Logout)
}

// SendClientToLogin 发送客户端登录请求。
// 该函数将指定的客户端添加到登录队列中，以处理登录操作。
// 参数:
//
//	client *Client: 待登录的客户端实例指针。
func (k *KafkaServer) SendClientToLogin(client *Client) {
	// 加锁以确保线程安全，防止多个goroutine同时修改Login通道。
	k.mutex.Lock()
	// 将客户端发送到Login通道，用于处理登录逻辑。
	k.Login <- client
	// 解锁以释放资源，允许其他goroutine执行相关操作。
	k.mutex.Unlock()
}

// SendClientToLogout 发送客户端登出请求。
// 该函数将指定的客户端添加到登出队列中，以处理登出操作。
// 参数:
//
//	client *Client: 待登出的客户端实例指针。
func (k *KafkaServer) SendClientToLogout(client *Client) {
	// 加锁以确保线程安全，防止多个goroutine同时修改Logout通道。
	k.mutex.Lock()
	// 将客户端发送到Logout通道，用于处理登出逻辑。
	k.Logout <- client
	// 解锁以释放资源，允许其他goroutine执行相关操作。
	k.mutex.Unlock()
}

// RemoveClient 从 KafkaServer 的 Clients 字典中移除指定的客户端。
// 参数:
//
//	uuid - 客户端的唯一标识符，用于定位并移除 Clients 字典中的相应条目。
//
// 该方法使用互斥锁来确保并发安全性，防止多个协程同时修改 Clients 字典。
func (k *KafkaServer) RemoveClient(uuid string) {
	// 加锁以确保接下来的操作是线程安全的。
	k.mutex.Lock()
	// 从 Clients 字典中删除指定 UUID 的客户端条目。
	delete(k.Clients, uuid)
	// 解锁以允许其他协程进行操作。
	k.mutex.Unlock()
}
