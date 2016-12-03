package ratelimiter

import (
	"fmt"
	"time"

	. "github.com/ctliu3/tap/store"
	"github.com/garyburd/redigo/redis"
)

type RedisLimiter struct {
	period time.Duration
	quota  int
	store  *RedisStore
}

type Pair struct {
	id    string
	value int
}

var Infinite = time.Duration(time.Hour * 24)

func NewRedisLimiter(storeOpt RedisStoreOption, period time.Duration, quota int) (*RedisLimiter, error) {
	store, err := NewRedisStore(storeOpt)
	if err != nil {
		return nil, err
	}

	return &RedisLimiter{
		period: period,
		quota:  quota,
		store:  store,
	}, nil
}

func (self *RedisLimiter) Acquire(key string) (time.Duration, error) {
	c := self.store.Pool.Get()
	defer c.Close()

	if err := c.Err(); err != nil {
		return Infinite, err
	}

	quota := &Pair{id: fmt.Sprintf("%s:%s:quota", self.store.Prefix, key)}
	remain := &Pair{id: fmt.Sprintf("%s:%s:remain", self.store.Prefix, key)}

	c.Send("WATCH", remain.id)
	defer c.Send("UNWATCH", remain.id)

	var (
		ret []interface{}
		err error
	)
	if ret, err = redis.Values(c.Do("MGET", quota.id, remain.id)); err != nil {
		return Infinite, err
	}

	if _, err = redis.Scan(ret, &quota.value, &remain.value); err != nil {
		return Infinite, err
	}

	if quota.value == 0 {
		c.Send("MULTI")
		c.Send("SET", quota.id, self.quota)
		c.Send("EXPIRE", quota.id, self.period.Seconds())
		c.Send("SET", remain.id, self.quota-1)
		c.Send("EXPIRE", remain.id, self.period.Seconds())
		if ret, err = redis.Values(c.Do("EXEC")); err != nil {
			return Infinite, err
		} else {
		}

		quota.value = self.quota
		remain.value = self.quota - 1
		return 0, nil
	} else if remain.value > 0 {
		c.Do("DECR", remain.id)
		remain.value--
		return 0, nil
	}

	return Infinite, fmt.Errorf("Everything is OK, but no quota available")
}
