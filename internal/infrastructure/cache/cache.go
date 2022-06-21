package cache

import (
	"context"
	"errors"
	"time"
)

var (
	// 过期时间
	DefaultExpiration time.Duration = 0
	ErrItemExpired    error         = errors.New("item has expired")
	ErrKeyNotFound    error         = errors.New("key not found in cache")
)

type Cache interface {
	Context(ctx context.Context) Cache
	Get(key string) (interface{}, time.Time, error)
	Put(key string, val interface{}, d time.Duration) error
	Delete(key string) error
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}

	return time.Now().UnixNano() > i.Expiration
}

func NewCache(opts ...Option) Cache {
	options := NewOptions(opts...)
	items := make(map[string]Item)

	if len(options.Items) > 0 {
		items = options.Items
	}

	return &memCache{
		opts:  options,
		items: items,
	}
}
