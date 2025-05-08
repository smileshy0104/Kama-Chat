package config

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	global.VIPER = core.Viper() // 加载配置文件（如 config.yaml）
	fmt.Println(global.CONFIG.MainConfig)
}
