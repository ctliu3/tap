package main

import (
	"fmt"
	"time"

	. "github.com/ctliu3/tap/ratelimiter"
)

func main() {
	limiter := NewTokenBucket(10, 10)

	for i := 0; i < 10; i++ {
		ok, dur := limiter.Acquire(400)
		if ok {
			fmt.Printf("scuc\n")
		} else {
			fmt.Printf("err, should wait %v\n", dur)
		}
		time.Sleep(time.Millisecond * 2)
	}
}
