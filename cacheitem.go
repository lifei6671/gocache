package gocache

import "time"

type Key string

type CacheItem struct {
	key          Key
	data         interface{}
	lifeDuration time.Duration
}
