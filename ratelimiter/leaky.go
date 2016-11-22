package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	sync.Mutex
	capacity   int
	perRequest time.Duration
	sleepFor   time.Duration
	maxSlack   time.Duration // for burstiness
	last       time.Time
}

// Parameter `rate' is request number the backend service can handle in each
// second, say, QPS.
// set capacity to 0 if you don't a buffer to handle burstiness
func NewLeakyBucket(opt BucketOption) (Bucket, error) {
	rate := opt.Rate
	capacity := opt.Capacity

	return &LeakyBucket{
		perRequest: time.Second / time.Duration(rate),
		capacity:   capacity,
		maxSlack:   -time.Duration(capacity) * time.Second / time.Duration(rate),
	}, nil
}

func (self *LeakyBucket) Name() string {
	return "LeakyBucket"
}

// maxWait: ms
func (self *LeakyBucket) Acquire(maxWait int) (bool, time.Duration) {
	self.Lock()
	defer self.Unlock()

	now := time.Now()

	if self.last.IsZero() {
		self.last = now
		return true, 0
	}

	waitDur := self.sleepFor + self.perRequest - now.Sub(self.last)
	if waitDur > time.Millisecond*time.Duration(maxWait) {
		return false, waitDur
	}

	self.sleepFor = waitDur
	// If sleepFor is negative, it means we have a buffer to handle the burstiness.
	// Too negative the sleepFor will lead to high overload if there are too many
	// requests in a short period of time.
	if self.sleepFor < self.maxSlack {
		self.sleepFor = self.maxSlack
	}

	if self.sleepFor > 0 {
		// fmt.Printf("dur:%v\n", self.sleepFor)
		time.Sleep(self.sleepFor)
		self.last = now.Add(self.sleepFor)
		self.sleepFor = 0
		return true, self.sleepFor
	}
	self.last = now
	return true, time.Duration(0)
}

func init() {
	Register("leaky", NewLeakyBucket)
}
