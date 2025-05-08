package main

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"fmt"
)

func main() {
	// 1. 配置管理初始化
	global.VIPER = core.Viper() // 加载配置文件（如 config.yaml）
	fmt.Println(global.CONFIG.MainConfig.Port)
}
