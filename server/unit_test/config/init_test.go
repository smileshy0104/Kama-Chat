package config

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	// 加载配置文件（如 config.yaml）
	global.VIPER = core.Viper()
	fmt.Println(global.CONFIG.MainConfig)
}
