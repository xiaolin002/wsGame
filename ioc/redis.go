package ioc

import "github.com/redis/go-redis/v9"

/**
 * @Description
 * @Date 2025/5/26 19:44
 **/

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "124.221.32.91:6379",
		Password: "123321", // no password set
	})
	return redisClient
}
