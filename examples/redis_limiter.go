package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ctliu3/tap/ratelimiter"
	"github.com/ctliu3/tap/store"
)

func main() {
	storeOpt := store.RedisStoreOption{
		Prefix:   "default",
		MaxConn:  10,
		Hostname: "127.0.0.1:6379",
	}

	limiter, err := ratelimiter.NewRedisLimiter(storeOpt, time.Duration(time.Second), 5)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	key := "ip1"
	size := 20
	done := make(chan bool, size)
	for i := 0; i < 20; i++ {

		go func(id int) {
			// Random in big range, so the request can hit in a long period (several
			// seconds).
			rnd := rand.Int() % 3000
			time.Sleep(time.Millisecond * time.Duration(rnd))
			dur, err := limiter.Acquire(key)
			if err != nil {
				fmt.Printf("id:%d, err: %v\n", id, err)
				done <- false
				return
			}

			if dur == 0 {
				fmt.Printf("id:%d, succ\n", id)
				done <- true
			} else {
				fmt.Printf("id:%d, fail, should wait for %v\n", id, dur)
				done <- false
			}
		}(i)

	}

	nOK := 0
	for i := 0; i < size; i++ {
		c := <-done
		if c {
			nOK += 1
		}
	}
	fmt.Printf("succ num: %v\n", nOK)
}
