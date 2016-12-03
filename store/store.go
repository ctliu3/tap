package store

type RedisStoreOption struct {
	Prefix   string
	MaxConn  int
	Hostname string
}
