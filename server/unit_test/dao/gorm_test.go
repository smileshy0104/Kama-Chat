package dao

import (
	"Kama-Chat/core"
	"Kama-Chat/global"
	"Kama-Chat/initialize/dao"
	"Kama-Chat/model"
	"Kama-Chat/utils/random"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	userInfo := &model.UserInfo{
		Uuid:      "U" + strconv.Itoa(random.GetRandomInt(11)),
		Nickname:  "apylee",
		Telephone: "13599185431",
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   1,
	}
	global.VIPER = core.Viper() // 加载配置文件（如 config.yaml）
	dao.InitMysql()
	_ = dao.GormDB.Create(userInfo)
}
