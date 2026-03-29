package cache

import (
	"sync"
	"time"
)

var expirations sync.Map

func SetWithExpiration(key, value string, d time.Duration) {
	Set(key, value)
	expTime := time.Now().Add(d)
	expirations.Store(key, expTime)
}

func ActiveExpiration() {
	go func() {
		for {
			time.Sleep(1 * time.Second)

			expirations.Range(func(k, v any) bool {
				key := k.(string)
				exp := v.(time.Time)

				if time.Now().After(exp) {
					Delete(key)
				}
				return true
			})
		}
	}()
}

func TTL(key string) time.Duration {
	_, exists := store.Load(key)
	if !exists {
		return -2
	}

	exp, ok := expirations.Load(key)
	if !ok {
		return -1
	}

	ttl := time.Until(exp.(time.Time))
	if ttl <= 0 {
		Delete(key)
		return -2
	}
	return ttl
}
