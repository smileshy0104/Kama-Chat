package random

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

// GetRandomInt 生成指定长度的随机整数。
// 参数 len 指定生成整数的位数。
// 返回值为生成的随机整数。
// 该函数通过计算确保生成的整数具有指定的位数。
func GetRandomInt(len int) int {
	return rand.Intn(9*int(math.Pow(10, float64(len-1)))) + int(math.Pow(10, float64(len-1)))
}

// GetNowAndLenRandomString 生成包含当前日期和指定长度随机整数的字符串。
// 参数 len 指定随机整数部分的位数。
// 返回值为格式化的当前日期和随机整数组成的字符串。
// 该函数首先获取当前日期并格式化为"20060102"的形式，然后调用 GetRandomInt 函数生成随机整数，
// 最后将两者拼接成一个字符串返回。
func GetNowAndLenRandomString(len int) string {
	return time.Now().Format("20060102") + strconv.Itoa(GetRandomInt(len))
}
