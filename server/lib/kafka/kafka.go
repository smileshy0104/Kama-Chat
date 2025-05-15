package kafka

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/zlog"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"time"
)

// ctx 用于在 Kafka 操作中传递上下文信息
var ctx = context.Background()

// kafkaService 结构体定义了 Kafka 相关的服务
type kafkaService struct {
	ChatWriter *kafka.Writer // ChatWriter 是用于发送消息的 Kafka Writer
	ChatReader *kafka.Reader // ChatReader 是用于接收消息的 Kafka Reader
	KafkaConn  *kafka.Conn   // KafkaConn 是用于创建 topic 的 Kafka 连接
	Ctx        *gin.Context  // Ctx 是 gin 上下文，用于在处理请求时传递上下文信息
}

// KafkaService 是 kafkaService 类型的全局实例
var KafkaService = new(kafkaService)

// KafkaInit 初始化 Kafka 服务。
// 该函数根据 global.CONFIG.KafkaConfig 中的配置设置 Kafka 的写入器和读取器。
// 不返回任何值。
func (k *kafkaService) KafkaInit() {
	// 获取全局配置中的 Kafka 配置
	kafkaConfig := global.CONFIG.KafkaConfig

	// 配置并初始化 Kafka 写入器，用于发送消息
	// 使用 Hash 负载均衡策略分配分区
	// 设置写入超时时间、无需确认机制，并禁用自动创建 Topic 功能
	k.ChatWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),   // Kafka Broker 地址
		Topic:                  kafkaConfig.ChatTopic,             // 发送消息的 Topic 名称
		Balancer:               &kafka.Hash{},                     // 使用 Hash 负载均衡策略分配分区
		WriteTimeout:           kafkaConfig.Timeout * time.Second, // 写入超时时间
		RequiredAcks:           kafka.RequireNone,                 // 不需要 broker 确认
		AllowAutoTopicCreation: false,                             // 禁止自动创建 Topic
	}

	// 配置并初始化 Kafka 读取器，用于接收消息
	// 设置 Broker 地址、监听的 Topic、提交间隔、消费者组 ID 和起始偏移量
	k.ChatReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},    // Kafka Broker 地址
		Topic:          kafkaConfig.ChatTopic,             // 消费的 Topic 名称
		CommitInterval: kafkaConfig.Timeout * time.Second, // 提交偏移量的时间间隔
		GroupID:        "chat",                            // 消费者组 ID
		StartOffset:    kafka.LastOffset,                  // 从最后一条消息开始消费
	})
}

// KafkaClose 关闭 Kafka 连接
// 该方法确保在服务停止时，正确关闭 ChatWriter 和 ChatReader
func (k *kafkaService) KafkaClose() {
	// 关闭 ChatWriter 和 ChatReader
	if err := k.ChatWriter.Close(); err != nil {
		zlog.Error(err.Error())
	}
	if err := k.ChatReader.Close(); err != nil {
		zlog.Error(err.Error())
	}
}

// CreateTopic 创建 Kafka topic
// 该方法检查是否存在指定的 topic，如果不存在则创建它
func (k *kafkaService) CreateTopic() {
	// 如果已经有topic了，就不创建了
	kafkaConfig := global.CONFIG.KafkaConfig

	// 获取全局配置中的 Kafka 配置ChatTopic
	chatTopic := kafkaConfig.ChatTopic

	// 连接至任意kafka节点
	var err error
	// 连接至kafka节点
	k.KafkaConn, err = kafka.Dial("tcp", kafkaConfig.HostPort)
	if err != nil {
		zlog.Error(err.Error())
	}

	// 创建topic的配置
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             chatTopic,             // topic名称
			NumPartitions:     kafkaConfig.Partition, // 分区数
			ReplicationFactor: 1,                     // 副本因子
		},
	}

	// 创建topic
	if err = k.KafkaConn.CreateTopics(topicConfigs...); err != nil {
		zlog.Error(err.Error())
	}
}
