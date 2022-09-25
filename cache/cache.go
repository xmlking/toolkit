package cache

import (
	"errors"

	lru "github.com/hashicorp/golang-lru"
)

type SetFn func() (v interface{}, err error)

type Cache interface {
	GetOrSet(k interface{}, setFn SetFn) (v interface{}, err error)
}

type GetSetCache struct {
	lru    *lru.Cache
	locker *ChanLocker
}

var (
	ErrCacheItemNotFound = errors.New("cache item not found")
)

func NewCache(size int) *GetSetCache {
	c, _ := lru.New(size)
	return &GetSetCache{
		lru:    c,
		locker: NewChanLocker(),
	}
}

func (c *GetSetCache) GetOrSet(k interface{}, setFn SetFn) (v interface{}, err error) {
	if val, ok := c.lru.Get(k); ok {
		return val, nil
	}
	acquired := c.locker.Lock(k, func() {
		v, err = setFn()
		if err != nil {
			return
		}
		c.lru.Add(k, v)
	})
	if acquired {
		return v, err
	}

	// someone else got the lock first and should have inserted something
	if v, ok := c.lru.Get(k); ok {
		return v, nil
	}

	// someone else acquired the lock, but no key was found
	// (most likely this value doesn't exist or the upstream fetch failed)
	return nil, ErrCacheItemNotFound
}
