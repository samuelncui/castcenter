package rediswrap

import (
	"strconv"
	"time"

	"github.com/abc950309/castcenter"
	"github.com/go-redis/redis"
)

// RedisWrap .
type RedisWrap struct {
	*redis.Client
}

var (
	_ = castcenter.Redis(&RedisWrap{})
)

// New .
func New(client *redis.Client) *RedisWrap {
	return &RedisWrap{Client: client}
}

// NewWithURL .
func NewWithURL(url string) (*RedisWrap, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return New(redis.NewClient(opt)), nil
}

// Get implement castcenter.Redis
func (w *RedisWrap) Get(key string) (string, error) {
	return w.Client.Get(key).Result()
}

// SetNX implement castcenter.Redis
func (w *RedisWrap) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return w.Client.SetNX(key, value, expiration).Result()
}

// Expire implement castcenter.Redis
func (w *RedisWrap) Expire(key string, expiration time.Duration) (bool, error) {
	return w.Client.Expire(key, expiration).Result()
}

// ZAdd implement castcenter.Redis
func (w *RedisWrap) ZAdd(key, member string, score float64) (int64, error) {
	return w.Client.ZAdd(key, &redis.Z{Member: member, Score: score}).Result()
}

// ZRangeByScore implement castcenter.Redis
func (w *RedisWrap) ZRangeByScore(key string, min, max float64) ([]string, error) {
	return w.Client.ZRangeByScore(key, &redis.ZRangeBy{
		Min: strconv.FormatFloat(min, 'f', -1, 64),
		Max: strconv.FormatFloat(max, 'f', -1, 64),
	}).Result()
}
