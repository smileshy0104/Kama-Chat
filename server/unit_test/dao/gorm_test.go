package dao

import (
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
		Telephone: "180323532112",
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   1,
	}
	_ = dao.GormDB.Create(userInfo)
}
