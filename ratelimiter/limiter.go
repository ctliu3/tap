package ratelimiter

import (
	"fmt"
	"log"
	"time"
)

type Limiter interface {
	Name() string
	Acquire(maxWait int) (bool, time.Duration)
}

type LimiterOption struct {
	Policy   string
	Rate     int
	Capacity int
}

type LimiterFactory func(opt LimiterOption) (Limiter, error)

// Store the mapping from the name to the class factory.
var limiterFactories = make(map[string]LimiterFactory)

func Register(policy string, factory LimiterFactory) {
	if factory == nil {
		log.Panicf("Limiter factory %s does not exists", policy)
	}
	_, registered := limiterFactories[policy]
	if registered {
		log.Fatalf("Limiter factory %s has registered", policy)
	}
	limiterFactories[policy] = factory
}

func CreateLimiter(opt LimiterOption) (Limiter, error) {
	paramLimiterPolicy := opt.Policy

	factory, ok := limiterFactories[paramLimiterPolicy]
	if !ok {
		return nil, fmt.Errorf("limiter policy %s has not registered", paramLimiterPolicy)
	}

	return factory(opt)
}
