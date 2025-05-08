package kafka

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/zlog"
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

// ctx 用于在 Kafka 操作中传递上下文信息
var ctx = context.Background()

// kafkaService 结构体定义了 Kafka 相关的服务
type kafkaService struct {
	ChatWriter *kafka.Writer
	ChatReader *kafka.Reader
	KafkaConn  *kafka.Conn
}

// KafkaService 是 kafkaService 类型的全局实例
var KafkaService = new(kafkaService)

// KafkaInit 初始化 Kafka 服务
// 该方法配置了 ChatWriter 和 ChatReader，为之后的消息发送和接收做准备
func (k *kafkaService) KafkaInit() {
	//k.CreateTopic()
	kafkaConfig := global.CONFIG.KafkaConfig
	k.ChatWriter = &kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.HostPort),
		Topic:                  kafkaConfig.ChatTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           kafkaConfig.Timeout * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.ChatReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaConfig.HostPort},
		Topic:          kafkaConfig.ChatTopic,
		CommitInterval: kafkaConfig.Timeout * time.Second,
		GroupID:        "chat",
		StartOffset:    kafka.LastOffset,
	})
}

// KafkaClose 关闭 Kafka 连接
// 该方法确保在服务停止时，正确关闭 ChatWriter 和 ChatReader
func (k *kafkaService) KafkaClose() {
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

	chatTopic := kafkaConfig.ChatTopic

	// 连接至任意kafka节点
	var err error
	k.KafkaConn, err = kafka.Dial("tcp", kafkaConfig.HostPort)
	if err != nil {
		zlog.Error(err.Error())
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             chatTopic,
			NumPartitions:     kafkaConfig.Partition,
			ReplicationFactor: 1,
		},
	}

	// 创建topic
	if err = k.KafkaConn.CreateTopics(topicConfigs...); err != nil {
		zlog.Error(err.Error())
	}

}
