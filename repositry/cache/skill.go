package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LockExpireTime = time.Second * 3
)

func SkillProductKey(productId uint) string {
	return fmt.Sprintf("sk:product:%d", productId)
}

func SkillLockKey(productId uint) string {
	return fmt.Sprintf("sk:lock:product:%d", productId)
}

type SkRedisCache interface {
	HMSetProductNumMoney(pId uint, num int, money float64) (bool, error)
	HMGetProductNumMoney(pId uint) (num int, money float64, err error)
	HDecrByNum(pId uint) (int64, error)

	Lock(pId uint, uuid string) (bool, error)
	UnLock(pId uint) (int64, error)
}

type customSkRedisCache struct {
	redisConn *redis.Client
}

func NewSkRedisCache(redisConn *redis.Client) SkRedisCache {
	return &customSkRedisCache{
		redisConn: redisConn,
	}
}

func (c *customSkRedisCache) HMSetProductNumMoney(pId uint, num int, money float64) (bool, error) {
	data := map[string]interface{}{
		"num":   num,
		"money": money,
	}
	return c.redisConn.HMSet(context.TODO(), SkillProductKey(pId), data).Result()
}

func (c *customSkRedisCache) HMGetProductNumMoney(pId uint) (num int, money float64, err error) {
	res, err := c.redisConn.HMGet(context.TODO(), SkillProductKey(pId), "num", "money").Result()
	if err != nil {
		return 0, 0, err
	}

	num, err = strconv.Atoi(res[0].(string))
	if err != nil {
		return 0, 0, err
	}

	money, err = strconv.ParseFloat(res[1].(string), 64)
	if err != nil {
		return 0, 0, err
	}

	return num, money, nil
}

func (c *customSkRedisCache) HDecrByNum(pId uint) (int64, error) {
	return c.redisConn.HIncrBy(context.TODO(), SkillProductKey(pId), "num", -1).Result()
}

func (c *customSkRedisCache) Lock(pId uint, uuid string) (bool, error) {
	return c.redisConn.SetNX(context.TODO(), SkillLockKey(pId), uuid, LockExpireTime).Result()
}

func (c *customSkRedisCache) UnLock(pId uint) (int64, error) {
	return c.redisConn.Del(context.TODO(), SkillLockKey(pId)).Result()
}
