package redisPool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-filestore-server/config"
	"time"
)

var (
	pool *redis.Pool
)

func InitRedis() {
	pool = newRedisPool()
	data, err := pool.Get().Do("KEYS", "*")
	fmt.Println(data, err)
}

func GetRedisPool() *redis.Pool {
	return pool
}

// 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			// 1.打开连接
			c, err := redis.Dial("tcp", config.DefaultConfig.RedisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// 2.访问认证
			if _, err = c.Do("AUTH", config.DefaultConfig.RedisPass); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
