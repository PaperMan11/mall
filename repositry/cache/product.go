package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	//RankKey 每日排名
	RankKey = "rank"
)

// 商品点击数缓存
func ProductViewKey(id uint) string {
	return fmt.Sprintf("view:product:%s", strconv.Itoa(int(id)))
}

type ProductRedisCache interface {
	ViewClick(productID uint) (uint64, error)
	AddClick(productID uint) error
}

type customProductRedisCache struct {
	redisConn *redis.Client
}

func NewProductRedisCache(redisConn *redis.Client) ProductRedisCache {
	return &customProductRedisCache{
		redisConn: redisConn,
	}
}

// View 获取点击数
func (c *customProductRedisCache) ViewClick(productID uint) (uint64, error) {
	s, err := c.redisConn.Get(context.TODO(), ProductViewKey(productID)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	count, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *customProductRedisCache) AddClick(productID uint) error {
	// 增加视频点击数
	_, err := c.redisConn.Incr(context.TODO(), ProductViewKey(productID)).Result()
	if err != nil {
		return err
	}
	// 增加排行点击数
	_, err = c.redisConn.ZIncrBy(context.TODO(), RankKey, 1, strconv.Itoa(int(productID))).Result()
	return err
}
