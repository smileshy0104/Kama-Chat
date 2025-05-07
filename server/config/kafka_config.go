package config

import "time"

type KafkaConfig struct {
	MessageMode string        `json:"message_mode" yaml:"message_mode"`
	HostPort    string        `json:"host_port" yaml:"host_port"`
	LoginTopic  string        `json:"login_topic" yaml:"login_topic"`
	LogoutTopic string        `json:"logout_topic" yaml:"logout_topic"`
	ChatTopic   string        `json:"chat_topic" yaml:"chat_topic"`
	Partition   int           `json:"partition" yaml:"partition"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
}
