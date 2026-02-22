package main

import (
	"sync"
	"time"
)

type Item struct {
	Value string
	Timer *time.Timer
}

var kv = struct {
	sync.RWMutex
	m map[string]*Item
}{m: make(map[string]*Item)}

func Set(key, value string) {
	kv.Lock()
	defer kv.Unlock()

	if item, exists := kv.m[key]; exists {
		if item.Timer != nil {
			item.Timer.Stop()
		}
	}

	kv.m[key] = &Item{Value: value}
}

func Get(key string) (string, bool) {
	kv.RLock()
	defer kv.RUnlock()

	val, ok := kv.m[key]
	if !ok {
		return "", false
	}
	return val.Value, ok
}

func Del(key string) bool {
	kv.Lock()
	defer kv.Unlock()

	item, exists := kv.m[key]
	if exists {
		if item.Timer != nil {
			item.Timer.Stop()
		}
		delete(kv.m, key)
		return true
	}
	return false
}

func Expire(key string, duration time.Duration) bool {
	kv.Lock()
	defer kv.Unlock()

	item, exists := kv.m[key]
	if !exists {
		return false
	}

	if item.Timer != nil {
		item.Timer.Stop()
	}

	item.Timer = time.AfterFunc(duration, func() {
		Del(key)
	})

	return true
}
