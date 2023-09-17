package tredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 使用示例
//rc := redis.ConnectRedis()
//defer rc.Close()
//redis.SetString(rc, "1236", "888", time.Minute)
//a := redis.GetString(rc, "1236")
//fmt.Println(a)

var ctx = context.Background()

// SetString 设置字符串型键值对
/*
@param redisClient *redis.Client
@param key string
@param value string
@param t time.Duration
*/
func SetString(redisClient *redis.Client, key string, value string, t time.Duration) {
	// 最后一个参数为保存时间，设为0的话就是永不过期
	_, err := redisClient.Set(ctx, key, value, t).Result()
	if err != nil {
		fmt.Println(err)
	}
}

// GetString 获取字符串型键值对
/*
@param redisClient *redis.Client
@param key string
@return value string
*/
func GetString(redisClient *redis.Client, key string) (value string) {
	value, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// key not exist
		value = "key not exist"
	} else if err != nil {
		// panic err
		value = "err"
	}
	return value
}

// DelByKey 删除键值对
func DelByKey(redisClient *redis.Client, key string) (int64, error) {
	result, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// GetExpiration 获取键值对的过期时间
func GetExpiration(redisClient *redis.Client, key string) (time.Duration, error) {
	expire, err := redisClient.TTL(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return expire, nil
}
