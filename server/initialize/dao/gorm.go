package dao

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormDB 是数据库的全局实例
var GormDB *gorm.DB

// init 函数在程序启动时初始化数据库连接
func init() {
	// 加载全局配置
	conf := global.CONFIG
	// 获取MySQL配置中的用户名
	user := conf.MysqlConfig.User
	// password := conf.MysqlConfig.Password
	// host := conf.MysqlConfig.Host
	// port := conf.MysqlConfig.Port
	// 获取主配置中的应用名称，用于数据库名称
	appName := conf.MainConfig.AppName
	// 构建MySQL的DSN（数据源名称）
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, appName)
	// 使用Unix套接字连接MySQL
	dsn := fmt.Sprintf("%s@unix(/var/run/mysqld/mysqld.sock)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, appName)
	// 初始化GormDB对象
	var err error
	// 使用Gorm库打开数据库连接
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 如果连接失败，记录错误日志并终止程序
	if err != nil {
		zlog.Fatal(err.Error())
	}
	// 自动迁移数据库模式，如果没有相应的表，会自动创建
	err = GormDB.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.UserContact{}, &model.Session{}, &model.ContactApply{}, &model.Message{})
	// 如果迁移失败，记录错误日志并终止程序
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
