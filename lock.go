package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

// DistributeLock 分布式锁结构
type DistributeLock struct {
	key        string             // key值
	expire     int                // 过期时间
	cancleFunc context.CancelFunc // 取消操作
	redis      *redis.Client      // redis客户端
}

// NewDistributeLock 创建分布式锁实例
func NewDistributeLock(rds *redis.Client, key string, expire int) *DistributeLock {
	return &DistributeLock{
		key:    key,
		expire: expire,
		redis:  rds,
	}
}

// Lock 加锁
func (l *DistributeLock) Lock() error {
	// 已锁定则返回err
	if err := l.lock(); err != nil {
		return err
	}
	// 创建一个上下文，若业务提前结束，需要关闭退出看门狗
	ctx, cancleFunc := context.WithCancel(context.Background())
	l.cancleFunc = cancleFunc
	l.startWatchDog(ctx)
	return nil
}

// 竞争redis锁
func (l *DistributeLock) lock() error {
	if err := l.redis.Do("SET", l.key, 1, "NX", "EX", l.expire).Err(); err != nil {
		return err
	}
	return nil
}

// 创建守护协程 自动续期
func (l *DistributeLock) startWatchDog(ctx context.Context) {
	go func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := l.redis.Do("EXPIRE", l.key, l.expire).Err(); err != nil {
					return err
				}
				time.Sleep(time.Duration(l.expire/2) * time.Second)
			}
		}
	}()
}
func (l *DistributeLock) Dtest() {

}

// Unlock 释放锁
func (l *DistributeLock) Unlock() error {
	if l.cancleFunc != nil {
		l.cancleFunc()
	}
	if err := l.redis.Do("Del", l.key).Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
