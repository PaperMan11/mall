package cache

import (
	"fmt"
	"mall/config"
	"testing"
)

func TestHash(t *testing.T) {
	c, err := NewRedisConn(&config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Db:       0,
		PoolSize: 100,
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	src := NewSkRedisCache(c)
	b, err := src.HMSetProductNumMoney(1, 10, 100.1)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(b)

	n, m, err := src.HMGetProductNumMoney(1)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(n, m)

	num, err := src.HDecrByNum(1)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(num)
}
