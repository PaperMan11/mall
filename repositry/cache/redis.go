package cache

import (
	"context"
	"fmt"
	"mall/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisConn(conf *config.RedisConf) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:       conf.Db,
		Password: conf.Password,
		PoolSize: conf.PoolSize,
	})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {

		return nil, err
	}
	return client, nil
}
