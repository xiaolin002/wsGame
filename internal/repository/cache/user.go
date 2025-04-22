package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	client redis.Cmdable
	// 这里可以添加缓存实现
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) GetUserFromCache(username string) (string, bool) {
	// 从缓存中获取用户信息
	return "", false // 示例，实际实现应返回缓存数据
}
