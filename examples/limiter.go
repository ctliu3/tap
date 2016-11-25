package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	. "github.com/ctliu3/tap/ratelimiter"
)

type Config struct {
	Option LimiterOption `toml:"limiter"`
}

func main() {
	confName := flag.String("conf", "conf.toml", "config file name")

	var config Config
	if _, err := toml.DecodeFile(*confName, &config); err != nil {
		fmt.Println(err)
		return
	}

	limiter, err := CreateLimiter(config.Option)
	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()
	for i := 0; i < 10; i++ {
		fmt.Printf("%v\n", time.Now().Sub(start))
		ok, dur := limiter.Acquire(1000)
		if ok {
			fmt.Printf("scuc\n")
		} else {
			fmt.Printf("err, should wait %v\n", dur)
		}
	}
}
