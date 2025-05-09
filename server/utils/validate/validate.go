package validate

import (
	"Kama-Chat/initialize/dao"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/model"
	"Kama-Chat/utils/constants"
	"errors"
	"gorm.io/gorm"
	"regexp"
)

// CheckTelephoneExist 检查手机号是否存在
func CheckTelephoneExist(telephone string) (string, int) {
	var user model.UserInfo
	// gorm默认排除软删除，所以翻译过来的select语句是SELECT * FROM `user_info` WHERE telephone = '18089596095' AND `user_info`.`deleted_at` IS NULL ORDER BY `user_info`.`id` LIMIT 1
	if res := dao.GormDB.Where("telephone = ?", telephone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("该电话不存在，可以注册")
			return "", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	zlog.Info("该电话已经存在，注册失败")
	return "该电话已经存在，注册失败", -2
}

// CheckTelephoneValid 检验电话是否有效
func CheckTelephoneValid(telephone string) bool {
	pattern := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	match, err := regexp.MatchString(pattern, telephone)
	if err != nil {
		zlog.Error(err.Error())
	}
	return match
}

// CheckEmailValid 校验邮箱是否有效
func CheckEmailValid(email string) bool {
	pattern := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		zlog.Error(err.Error())
	}
	return match
}

// CheckUserIsAdminOrNot 检验用户是否为管理员
func CheckUserIsAdminOrNot(user model.UserInfo) int8 {
	return user.IsAdmin
}
