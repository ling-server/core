package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ling-server/core/log"
)

var (
	fetchOrSaveMutex = keyMutex{m: &sync.Map{}}
)

// FetchOrSave retrieves the value for the key if present in the cache.
// Otherwise, it saves the value from the builder and retrieves the value
// for the key again.
func FetchOrSave(ctx context.Context, c Cache, key string, value interface{},
	builder func() (interface{}, error), expiration ...time.Duration) error {
	err := c.Fetch(ctx, key, value)
	// value found from the cache
	if err == nil {
		return nil
	}

	// internal error
	if !errors.Is(err, ErrorNotFound) {
		return err
	}

	// lock the key in cache and try to build the value for the key
	lockKey := fmt.Sprintf("%p:%s", c, key)
	fetchOrSaveMutex.Lock(lockKey)
	defer fetchOrSaveMutex.Unlock(lockKey)

	// fetch agai to avoid to build the value multi-times
	err = c.Fetch(ctx, key, value)
	if err == nil {
		return nil
	}

	// internal error
	if !errors.Is(err, ErrorNotFound) {
		return err
	}

	val, err := builder()
	if err != nil {
		return err
	}

	if err := c.Save(ctx, key, val, expiration...); err != nil {
		log.Warningf("Failed to save value to cache, error: %v", err)

		// save the value to cache failed, copy it to the value using the default
		// codec
		data, err := codec.Encode(val)
		if err != nil {
			return err
		}

		return codec.Decode(data, value)
	}

	// after the building, fetch value again
	return c.Fetch(ctx, key, val)
}
