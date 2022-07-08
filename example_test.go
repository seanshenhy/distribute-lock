package lock_test

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	lock "github.com/seanshenhy/distribute-lock/v2"
)

// 测试锁
func Example() {
	key := "test"
	rds := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	dl := lock.NewDistributeLock(rds, key, 2)
	if err := dl.Lock(); err != nil {
		panic(err)
	}
	go func() {
		for {
			res, err := rds.Get(key).Result()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(key, res)
			}
			// Output: test 1
			// test 1
			// test 1
			// test 1
			time.Sleep(time.Duration(1 * time.Second))
		}
	}()
	doJob()
	dl.Unlock()
}
func doJob() {
	time.Sleep(4 * time.Second)
}
