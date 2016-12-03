package store

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type RedisStore struct {
	Prefix  string
	Pool    *redis.Pool
	MaxConn int
}

func NewRedisStore(opt RedisStoreOption) (*RedisStore, error) {
	pool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", opt.Hostname)
		if err != nil {
			return nil, err
		}
		return c, err
	}, opt.MaxConn)

	if pool == nil {
		return nil, fmt.Errorf("Init redis pool error")
	}

	store := &RedisStore{
		Prefix:  opt.Prefix,
		Pool:    pool,
		MaxConn: opt.MaxConn,
	}

	_, err := store.ping() // test connection
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (self *RedisStore) ping() (bool, error) {
	conn := self.Pool.Get()
	defer conn.Close()

	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}
