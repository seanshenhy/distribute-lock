# 分布式锁0.0.1版本
这是一个玩具分布式锁，它是基于续期实现，每次加锁成功都会创建一个看门狗，释放锁会移出看门狗。

# 快速使用

安装
```shell
go get  github.com/seanshenhy/distribute-lock
```
使用

```shell
key := "test"
rds := redis.NewClient(&redis.Options{
    Addr: ":6379",
})
dl := lock.NewDistributeLock(rds, key, 2)
if err := dl.Lock(); err != nil {
    panic(err)
}
func (){
    // do job
    time.Sleep(time.Duration(5 * time.Second))
}()
dl.Unlock()
```
# 后续
后面会逐渐改造变为工业级使用，欢迎交流。