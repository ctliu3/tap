package ratelimiter

import (
	"fmt"
	"log"
	"time"
)

type Bucket interface {
	Name() string
	Acquire(maxWait int) (bool, time.Duration)
}

type BucketOption struct {
	Name     string
	Rate     int
	Capacity int
}

type BucketFactory func(opt BucketOption) (Bucket, error)

// Store the mapping from the name to the class factory.
var bucketFactories = make(map[string]BucketFactory)

func Register(name string, factory BucketFactory) {
	if factory == nil {
		log.Panicf("Bucket factory %s does not exists", name)
	}
	_, registered := bucketFactories[name]
	if registered {
		log.Fatalf("Bucket factory %s has registered", name)
	}
	bucketFactories[name] = factory
}

func CreateBucket(opt BucketOption) (Bucket, error) {
	paramBucketName := opt.Name

	factory, ok := bucketFactories[paramBucketName]
	if !ok {
		return nil, fmt.Errorf("bucket name %s has not registered", paramBucketName)
	}

	return factory(opt)
}
