package main

import (
	"fmt"
	"time"

	. "github.com/ctliu3/tap/ratelimiter"
)

func main() {
	limiter := NewLeakyBucket(100)

	prev := time.Now()
	for i := 0; i < 10; i++ {
		cur := limiter.Attempt()
		fmt.Printf("%v\n", cur.Sub(prev))
		prev = cur
	}
}
