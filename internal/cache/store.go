package cache

import "sync"

// sync map over normal map due to concurrent access via goroutines (can also use map + mutex)
var store sync.Map

func Set(key, value string) {
	store.Store(key, value)
}

func Get(key string) (string, bool) {
	val, ok := store.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}

func Delete(key string) bool {
	_, ok := store.LoadAndDelete(key)
	return ok
}

func LoadAndDelete(key string) (string, bool) {
	val, ok := store.LoadAndDelete(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}
