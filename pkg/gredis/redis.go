package gredis

import (
	"encoding/json"
	"gin-blog/pkg/setting"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisCoon *redis.Pool

func Setup() error {
	RedisCoon = &redis.Pool{
		//最大空闲连接数
		MaxIdle: setting.RedisSetting.MaxIdle,
		//在给定时间内，允许分配的最大连接数（当为零时，没有限制）
		MaxActive: setting.RedisSetting.MaxActive,
		//在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）
		IdleTimeout:     setting.RedisSetting.IdleTimeout,
		Wait:            false,
		MaxConnLifetime: 0,
		//提供创建和配置应用程序连接的一个函数
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		//可选的应用程序检查健康功能
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
	conn := RedisCoon.Get()
	defer conn.Close()
	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	reply, err := redis.Bool(conn.Do("SET", key, value))
	_, _ = conn.Do("EXPIRE", key, time)
	return reply, err
}

func Exists(key string) bool {
	conn := RedisCoon.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisCoon.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, err
}

func Delete(key string) (bool, error) {
	conn := RedisCoon.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisCoon.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
