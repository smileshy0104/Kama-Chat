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
	// 获取 Redis 配置
	conf := global.CONFIG
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.RedisConfig.Db
	addr := host + ":" + strconv.Itoa(port)

	// 创建 Redis 客户端实例
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试 Redis 连接
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		zlog.Error(fmt.Sprintf("Redis连接失败,err:%v", err))
	} else {
		zlog.Info(fmt.Sprintf("Redis连接成功,pong:%v", pong))
	}

}

// SetKeyEx 设置Redis键的值，并指定过期时间。
// 该函数用于在Redis中存储一个键值对，并且可以为该键设置一个过期时间。
func SetKeyEx(key string, value string, timeout time.Duration) error {
	// 使用redisClient.Set方法设置Redis键的值和过期时间，并检查操作是否成功。
	err := redisClient.Set(ctx, key, value, timeout).Err()
	// 如果操作失败，直接返回错误信息。
	if err != nil {
		return err
	}
	// 如果操作成功，返回nil表示操作完成且无错误。
	return nil
}

// GetKey 从Redis中获取指定key的值。
func GetKey(key string) (string, error) {
	// 使用redisClient.Get方法从Redis中获取键对应的值。
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		// 当获取的值发生错误时，检查错误类型。
		if errors.Is(err, redis.Nil) {
			// 如果错误类型为redis.Nil，表示key不存在。
			zlog.Info("该key不存在")
			return "", nil
		}
		// 如果错误类型不是redis.Nil，则返回错误信息。
		return "", err
	}
	// 如果没有发生错误，则返回获取到的值。
	return value, nil
}

// GetKeyNilIsErr 从Redis中获取指定键的值。
// 如果遇到任何错误（包括键不存在的情况），都将返回一个错误。
func GetKeyNilIsErr(key string) (string, error) {
	// 使用redisClient.Get方法从Redis中获取键对应的值。
	// 如果操作失败或键不存在，err将不为nil。
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		// 如果发生错误，返回空字符串和错误详情。
		return "", err
	}
	// 如果操作成功且键存在，返回键对应的值和nil作为错误。
	return value, nil
}

// GetKeyWithPrefixNilIsErr 根据前缀获取键，如果未找到键或找到多个键则返回错误。
func GetKeyWithPrefixNilIsErr(prefix string) (string, error) {
	// 定义一个变量用于存储键名
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, prefix+"*").Result()
		if err != nil {
			return "", err
		}

		// 如果没有找到任何键，返回空字符串和redis.Nil错误
		if len(keys) == 0 {
			zlog.Info("没有找到相关前缀key")
			return "", redis.Nil
		}

		// 如果找到唯一的键，返回该键
		if len(keys) == 1 {
			zlog.Info(fmt.Sprintln("成功找到了相关前缀key", keys))
			return keys[0], nil
		} else {
			// 如果找到多个键，记录错误并返回自定义错误
			zlog.Error("找到了数量大于1的key，查找异常")
			return "", errors.New("找到了数量大于1的key，查找异常")
		}
	}
}

// GetKeyWithSuffixNilIsErr 根据给定的后缀查找匹配的Redis键。
// 如果找到一个匹配的键，则返回该键；如果未找到任何键，则返回空字符串和redis.Nil错误。
// 如果找到多个键，视为异常情况，返回错误。
func GetKeyWithSuffixNilIsErr(suffix string) (string, error) {
	// 定义一个变量用于存储键名
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
func DelKeyIfExists(key string) error {
	// 使用 Exists 方法检查键是否存在
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
			// 使用 Del 命令删除键
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
func DelKeysWithPrefix(prefix string) error {
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
			// 使用 Del 命令删除键
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
func DelKeysWithSuffix(suffix string) error {
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

// DeleteAllRedisKeys 删除Redis中所有的键。
// 该函数通过循环扫描Redis键空间，并分批删除找到的键，以减少操作对Redis性能的影响。
// 注意：此操作可能会影响性能，谨慎在生产环境中使用。
func DeleteAllRedisKeys() error {
	// 初始化游标为0，用于开始扫描。
	var cursor uint64 = 0
	for {
		// 使用Scan命令分批获取键，"*"表示匹配所有键，0表示每次扫描返回的键的数量无限制。
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			// 如果扫描过程中出现错误，返回错误。
			return err
		}
		// 更新游标为下一次扫描的起始点。
		cursor = nextCursor

		// 如果扫描到的键数量大于0，则尝试删除这些键。
		if len(keys) > 0 {
			// 使用Del命令删除之前扫描得到的键。
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				// 如果删除过程中出现错误，返回错误。
				return err
			}
		}

		// 如果游标回到0，表示所有键都已经扫描并处理完毕，退出循环。
		if cursor == 0 {
			break
		}
	}
	// 所有键删除成功，返回nil表示操作成功。
	return nil
}
