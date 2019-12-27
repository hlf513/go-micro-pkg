package etcd

import (
	"github.com/micro/go-micro/util/log"
	"go.etcd.io/etcd/clientv3"
)

// Client 创建 etcd client
func Connect() (*clientv3.Client, error) {
	etcdConf, err := GetEtcdConf()
	if err != nil {
		log.Fatal("[etcd connect] ", err.Error())
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints: etcdConf.Host,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 分布式锁
// func Locker(tryLockTimeout ...time.Duration) (*etcd.Client, lock.Locker, error) {
// 	tlt := 1 * time.Second
// 	if len(tryLockTimeout) != 0 {
// 		tlt = tryLockTimeout[0]
// 	}
// 	client, err := Client()
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	locker := lock.NewEtcdLocker(client, lock.WithTrylockTimeout(tlt))
// 	return client, locker, nil
// }
// 
// // usage:
// func example() {
// 	// 初始化 etcdLocker（默认 trylocktiemout 1s，可传参自定义）
// 	client, locker, err := Locker()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	defer client.Close()
// 	// 加锁(tryLock 模式)，超时时间使用 trylocktimeout
// 	lock, err := locker.Acquire("/lock", 10)
// 	// lock, err := locker.WaitAcquire("/lock", 10) // WaitAcquire 阻塞模式，直到获取锁
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	// 这里写逻辑，但是不要超过自定义的 10s，超过10s后，其他进程可以直接获取锁
// 
// 	// 释放锁 
// 	err = lock.Release()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// }
