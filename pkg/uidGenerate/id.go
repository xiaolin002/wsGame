package uidGenerate

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

// RedisIncrClient 定义 Redis 操作接口
type RedisIncrClient interface {
	Incr(ctx context.Context, key string) (int64, error)
}

type UIDGenerator struct {
	serverID    uint64 // 区服ID
	roleType    uint64 // 角色类型
	sequenceMax uint64 // 自增序列最大值（根据位数计算）
	mutex       sync.Mutex
	redisClient redis.Cmdable // 实现 RedisIncrClient 接口的客户端
}

// NewUIDGenerator 初始化UID生成器
// serverID: 区服ID（0~99999）
// roleType: 角色类型（0~99）
// redisClient: 实现 RedisIncrClient 接口的客户端
func NewUIDGenerator(serverID, roleType uint64, redisClient redis.Cmdable) (*UIDGenerator, error) {
	// 校验参数范围
	if serverID > 99999 {
		return nil, errors.New("serverID must be <= 99999")
	}
	if roleType > 99 {
		return nil, errors.New("roleType must be <= 99")
	}

	return &UIDGenerator{
		serverID:    serverID,
		roleType:    roleType,
		sequenceMax: 999999, // 6位自增序列
		redisClient: redisClient,
	}, nil
}

// Generate 生成唯一UID（线程安全）
func (g *UIDGenerator) Generate(ctx context.Context) (uint64, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// 获取当前日期
	currentDate := time.Now().Format("20060102")
	// 以日期作为 Redis 键的一部分
	redisKey := "uid_counter:" + currentDate

	// 使用 Redis 的 INCR 命令自增 UID
	cmd := g.redisClient.Incr(ctx, redisKey)
	newUID, err := cmd.Result()
	if err != nil {
		return 0, err
	}

	// 转换 newUID 为 uint64 类型
	newUIDUint64 := uint64(newUID)

	// 检查自增序列是否溢出
	if newUIDUint64 > g.sequenceMax {
		return 0, errors.New("sequence overflow")
	}

	// 将日期转换为数字
	dateNum, err := strconv.ParseUint(currentDate, 10, 64)
	if err != nil {
		return 0, err
	}

	// 计算各字段的位数偏移
	uid := g.serverID * 1e13 // 5位区服ID → 占前5位
	uid += g.roleType * 1e11 // 2位角色 → 接5位后
	uid += dateNum * 1e6     // 8位日期 → 接7位后
	uid += newUIDUint64      // 使用新的 UID

	return uid, nil
}
