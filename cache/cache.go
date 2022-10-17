package cache

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ling-server/core/log"
	"github.com/ling-server/core/retry"
)

const (
	// Memory the cache name of memory
	Memory = "memory"
	// Redis the cache name of redis
	Redis = "redis"
	// RedisSentinel the cache name of redis sentinel
	RedisSentinel = "redis+sentinel"
)

var (
	// ErrorNotFound error returns the key value not found in the cache
	ErrorNotFound = errors.New("Key not found")
)

type Cache interface {
	Contain(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
	Fetch(ctx context.Context, key string, value interface{}) error
	Ping(ctx context.Context) error
	Save(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error
	Keys(ctx context.Context, prefixes ...string) ([]string, error)
}

var (
	factories      = map[string]func(opts Options) (Cache, error){}
	factoriesMutex sync.RWMutex
)

func Register(t string, factory func(opts Options) (Cache, error)) {
	factoriesMutex.Lock()
	defer factoriesMutex.Unlock()

	factories[t] = factory
}

func New(t string, opt ...Option) (Cache, error) {
	opts := newOptions(opt...)
	opts.Codec = codec // use the default codec for the cache

	factoriesMutex.Lock()
	defer factoriesMutex.Unlock()

	factory, ok := factories[t]
	if !ok {
		return nil, fmt.Errorf("Cache type %s not support", t)
	}

	return factory(opts)
}

var (
	cache Cache
)

func Initialize(t, addr string) error {
	c, err := New(t, Address(addr), Prefix("cache:"))
	if err != nil {
		return err
	}

	redactedAddr := addr
	if u, err := url.Parse(addr); err == nil {
		redactedAddr = redacted(u)
	}

	opts := []retry.Option{
		retry.InitialInterval(time.Millisecond * 500),
		retry.MaxInterval(time.Second * 10),
		retry.Timeout(time.Minute),
		retry.Callback(func(err error, sleep time.Duration) {
			log.Errorf("Failed to ping %s, retry after %s : %v", redactedAddr, sleep, err)
		}),
	}

	if err := retry.Retry(func() error {
		return c.Ping(context.TODO())
	}, opts...); err != nil {
		return err
	}

	cache = c
	return nil
}

func Default() Cache {
	return cache
}
