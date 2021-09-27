package main

import (
	"context"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)

func main() {
	// setAndGet()
	mutex()
}

func setAndGet() {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{"192.168.100.133:1000"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "q1mi", "dsb")
	// 这里cli不管成没成功也已同步执行完了 所以可以立即调用cancel释放操作进程
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "q1mi")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

func mutex() {
	// cli, err := etcd.New(etcd.Config{
	// 	Endpoints:   []string{"192.168.100.133:1000"},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if err != nil {
	// 	fmt.Printf("connect to etcd failed, err:%v\n", err)
	// 	return
	// }
	// fmt.Println("connect to etcd success")
	// defer cli.Close()
	// 创建两个人会话
}
