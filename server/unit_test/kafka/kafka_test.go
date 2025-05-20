package kafka

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/lib/chat"
	"Kama-Chat/lib/kafka"
	myredis "Kama-Chat/lib/redis"
	"testing"
)

func TestKafka(t *testing.T) {
	// 1. 配置管理初始化
	global.VIPER = core.Viper() // 加载配置文件（如 config.yaml）
	// 2. 日志初始化
	zlog.InitLogger()
	// 3. 数据库初始化
	dao.InitMysql()
	// 4. Redis初始化
	myredis.InitRedis()

	// 5. kafka初始化
	chat.InitKafka()
	kafka.KafkaService.KafkaInit()
	defer kafka.KafkaService.KafkaClose()
	kafka.KafkaService.CreateTopic()
}
