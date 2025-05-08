package core

import (
	"Kama-Chat/core/internal"
	"Kama-Chat/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper 初始化 viper 配置实例。
// 参数 path 是可选的配置文件路径，如果没有提供，则使用默认路径。
// 返回值是初始化后的 *viper.Viper 实例。
func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		config = internal.ConfigDefaultFile // 默认值 "config.yaml"
	} else {
		config = path[0]
	}
	v := viper.New()
	v.SetConfigFile(config) // 👈 设置完整路径+文件名
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.CONFIG); err != nil {
		panic(err)
	}

	return v
}
