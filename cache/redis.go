package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"simple-web-golang/config"
	"time"
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Exists(key string) (bool, error)
	UpdateRank(key string, point int, uid int64) error
	Rank(key string, uid int64) (int64, error)
	Flush() error
	SetMap(key string, value map[string]interface{}) error
	GetMap(key string) (map[string]interface{}, error)
}

func FromContext(c context.Context) Cache {
	val := c.Value(config.KeyCache)
	return val.(Cache)
}

type Redis struct {
	Cli  *redis.Client
	Conf *config.CacheConfig
}

func New(c *config.CacheConfig) *Redis {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", c.Host, c.Port),
		DB:   c.Db,
	})
	if _, err := cli.Ping().Result(); err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}

	return &Redis{
		Cli:  cli,
		Conf: c,
	}
}

func (r *Redis) prefixed(key string) string {
	return fmt.Sprint(r.Conf.Prefix, key)
}

func (r *Redis) Set(key string, value interface{}) error {
	key = r.prefixed(key)
	return r.Cli.Set(key, value, time.Duration(r.Conf.Expire)*time.Second).Err()
}
func (r *Redis) Get(key string) (interface{}, error) {
	key = r.prefixed(key)
	var ret interface{}
	err := r.Cli.Get(key).Scan(ret)
	return ret, err
}
func (r *Redis) Exists(key string) (bool, error) {
	key = r.prefixed(key)
	val, err := r.Cli.Exists(key).Result()
	return val != int64(0), err
}
func (r *Redis) UpdateRank(key string, point int, uid int64) error {
	key = r.prefixed(key)
	return r.Cli.ZAdd(key, redis.Z{
		Score:  float64(point),
		Member: uid,
	}).Err()
}
func (r *Redis) Rank(key string, uid int64) (int64, error) {
	key = r.prefixed(key)
	return r.Cli.ZRank(key, string(uid)).Result()
}
func (r *Redis) Flush() error {
	return r.Cli.FlushDB().Err()
}
func (r *Redis) SetMap(key string, value map[string]interface{}) error {
	key = r.prefixed(key)
	return r.Cli.HMSet(key, value).Err()
}
func (r *Redis) GetMap(key string) (map[string]interface{}, error) {
	key = r.prefixed(key)
	ret := make(map[string]interface{})
	get, err := r.Cli.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, config.ErrCacheMissingKey
	}
	for key, val := range get {
		ret[key] = val
	}
	return ret, nil
}
