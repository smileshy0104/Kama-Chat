package sms

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/zlog"
	"Kama-Chat/lib/redis"
	"Kama-Chat/utils/constants"
	"Kama-Chat/utils/random"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strconv"
	"time"
)

var smsClient *dysmsapi20170525.Client

// createClient 使用AK&SK初始化账号Client
func createClient() (result *dysmsapi20170525.Client, err error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	accessKeyID := global.CONFIG.AuthCodeConfig.AccessKeyID
	accessKeySecret := global.CONFIG.AuthCodeConfig.AccessKeySecret
	if smsClient == nil {
		config := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(accessKeyID),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(accessKeySecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
		smsClient, err = dysmsapi20170525.NewClient(config)
	}
	return smsClient, err
}

// VerificationCode 向指定电话号码发送验证码，并返回发送结果和错误代码
// 该函数首先尝试从Redis中获取未过期的验证码，如果存在，则直接返回提示信息
// 如果验证码已过期或不存在，则生成新的验证码，存储到Redis，并发送到指定的电话号码
func VerificationCode(telephone string) (string, int) {
	// 创建客户端
	client, err := createClient()
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	// 构造Redis键名
	key := "auth_code_" + telephone

	// 尝试从Redis获取验证码
	code, err := redis.GetKey(key)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	// 如果验证码存在且未过期
	if code != "" {
		// 直接返回，验证码还没过期，用户应该去输验证码
		message := "目前还不能发送验证码，请输入已发送的验证码"
		zlog.Info(message)
		return message, -2
	}

	// 验证码过期，重新生成
	code = strconv.Itoa(random.GetRandomInt(6))
	fmt.Println(code)

	// 将新的验证码存储到Redis，并设置过期时间
	err = redis.SetKeyEx(key, code, time.Hour*24) // 1分钟有效
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	// 准备发送短信的请求参数
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"), // 短信模板
		PhoneNumbers:  tea.String(telephone),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
	}

	// 创建运行时选项
	runtime := &util.RuntimeOptions{}
	// 目前使用的是测试专用签名，签名必须是“阿里云短信测试”，模板code为“SMS_154950909”
	rsp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}

	// 记录发送短信的响应
	zlog.Info(*util.ToJSONString(rsp))
	// 返回成功提示信息
	return "验证码发送成功，请及时在对应电话查收短信", 0
}
