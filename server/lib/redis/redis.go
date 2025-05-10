package redis

import (
	"Kama-Chat/global"
	"Kama-Chat/initialize/zlog"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

// redisClient 是 Redis 客户端实例
var redisClient *redis.Client

// ctx 是用于 Redis 操作的上下文
var ctx = context.Background()

// init 函数在包被导入时初始化 Redis 客户端
func InitRedis() {
	conf := global.CONFIG
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.RedisConfig.Db
	addr := host + ":" + strconv.Itoa(port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// SetKeyEx 设置 Redis 中的键值对，并指定过期时间
// 参数:
//
//	key: 键名
//	value: 键值
//	timeout: 键的过期时间
//
// 返回值:
//
//	error: 错误信息，如果设置成功则为nil
func SetKeyEx(key string, value string, timeout time.Duration) error {
	err := redisClient.Set(ctx, key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetKey 从 Redis 中获取键的值
// 参数:
//
//	key: 键名
//
// 返回值:
//
//	string: 键的值，如果键不存在则为空字符串
//	error: 错误信息，如果获取成功则为nil
func GetKey(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			zlog.Info("该key不存在")
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// GetKeyNilIsErr 从 Redis 中获取键的值，如果键不存在则返回错误
// 参数:
//
//	key: 键名
//
// 返回值:
//
//	string: 键的值
//	error: 错误信息，如果键不存在则为redis.Nil
func GetKeyNilIsErr(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetKeyWithPrefixNilIsErr 获取 Redis 中带有指定前缀的键，如果不存在则返回错误
// 参数:
//
//	prefix: 键名前缀
//
// 返回值:
//
//	string: 匹配的键名
//	error: 错误信息，如果没有找到匹配的键则为redis.Nil
func GetKeyWithPrefixNilIsErr(prefix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return "", err
		}

		if len(keys) == 0 {
			zlog.Info("没有找到相关前缀key")
			return "", redis.Nil
		}

		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关前缀key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于1的key，查找异常")
			return "", errors.New("找到了数量大于1的key，查找异常")
		}
	}

}

// GetKeyWithSuffixNilIsErr 获取 Redis 中带有指定后缀的键，如果不存在则返回错误
// 参数:
//
//	suffix: 键名后缀
//
// 返回值:
//
//	string: 匹配的键名
//	error: 错误信息，如果没有找到匹配的键则为redis.Nil
func GetKeyWithSuffixNilIsErr(suffix string) (string, error) {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
		if err != nil {
			return "", err
		}

		if len(keys) == 0 {
			zlog.Info("没有找到相关后缀key")
			return "", redis.Nil
		}

		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关后缀key", keys))
			return keys[0], nil
		} else {
			zlog.Error("找到了数量大于1的key，查找异常")
			return "", errors.New("找到了数量大于1的key，查找异常")
		}
	}

}

// DelKeyIfExists 如果键存在，则删除该键
// 参数:
//
//	key: 键名
//
// 返回值:
//
//	error: 错误信息，如果键不存在则为nil
func DelKeyIfExists(key string) error {
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 1 { // 键存在
		delErr := redisClient.Del(ctx, key).Err()
		if delErr != nil {
			return delErr
		}
	}
	// 无论键是否存在，都不返回错误
	return nil
}

// DelKeysWithPattern 删除 Redis 中匹配指定模式的键
// 参数:
//
//	pattern: 键名模式
//
// 返回值:
//
//	error: 错误信息，如果操作成功则为nil
func DelKeysWithPattern(pattern string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			log.Println("没有找到对应key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相关对应key", keys)
		}
	}

	return nil
}

// DelKeysWithPrefix 删除 Redis 中带有指定前缀的键
// 参数:
//
//	prefix: 键名前缀
//
// 返回值:
//
//	error: 错误信息，如果操作成功则为nil
func DelKeysWithPrefix(prefix string) error {
	//var cursor uint64 = 0
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			log.Println("没有找到相关前缀key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相关前缀key", keys)
		}
	}

	return nil
}

// DelKeysWithSuffix 删除 Redis 中带有指定后缀的键
// 参数:
//
//	suffix: 键名后缀
//
// 返回值:
//
//	error: 错误信息，如果操作成功则为nil
func DelKeysWithSuffix(suffix string) error {
	//var cursor uint64 = 0
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, "*"+suffix).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			log.Println("没有找到相关后缀key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相关后缀key", keys)
		}
	}

	return nil
}

// DeleteAllRedisKeys 删除 Redis 中的所有键
// 返回值:
//
//	error: 错误信息，如果操作成功则为nil
func DeleteAllRedisKeys() error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		cursor = nextCursor

		if len(keys) > 0 {
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
