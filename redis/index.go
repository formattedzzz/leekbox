package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
)

var DB redis.Client
var ctx = context.Background()

func init() {
	DB = *redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})
	fmt.Println("DB connected.")
}
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func main() {
	// T0()
	// T1()
	T2()
	// T3()
}

func T0() {
	if res, err := DB.HGet(ctx, "hash", "info").Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(100))
	fmt.Println(RandString(20))
	if res, err := DB.Exists(ctx, "project").Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	if res, err := DB.LInsertAfter(ctx, "array", "npm", "inserted").Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	if res, err := DB.Set(ctx, "leo", 123, time.Minute).Result(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res == "OK")
	}
	if num, err := DB.Get(ctx, "leo").Int(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num, num+2)
	}
}

func T1() {
	cb := func(tx *redis.Tx) error {
		fmt.Println("watched")
		_, err := tx.Set(ctx, "leo", 123, 0).Result()
		if err != nil {
			return err
		}
		num, err := DB.Get(ctx, "leo").Int()
		if err != nil {
			return err
		}
		num += 2
		fmt.Println("num", num)
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if _, err := pipe.Set(ctx, "leo2", num, 0).Result(); err != nil {
				return err
			}
			return nil
		})
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if _, err := pipe.Set(ctx, "leo3", num+2, 0).Result(); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
	DB.Watch(ctx, cb, "project")
}

func T2() {
	// pipe := DB.TxPipeline()
	// pipe.Set(ctx, "leo", "caoya", 0)
	// pipe.Set(ctx, "leo2", "caoya", 0)
	// pipe.SetNX(ctx, "count", 123, 0)
	// pipe.IncrBy(ctx, "count", 10)
	// intcmd := pipe.Get(ctx, "count").Val()
	// pipe.Exec(ctx)
	// fmt.Println(intcmd)

	cmds, err := DB.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.Append(ctx, "leo", "ooo")
		p.IncrBy(ctx, "count", 10)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cmds[0].(*redis.IntCmd).Val(), cmds[1].(*redis.IntCmd).Val())
}

func T3() {
	// 这个监听主要是负责多个客户端连接同一个redis-server的锁键问题
	// 如果一个有另一个客户端已经在监听了 那么这个监听的事务在exec的时候会失败
	DB.Watch(ctx, func(t *redis.Tx) error {
		fmt.Println("触发监听")
		t.Incr(ctx, "count")
		return nil
	}, "leo")
	signChan := make(chan os.Signal)
	signal.Notify(signChan, os.Interrupt)
	<-signChan
}
