package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data           map[string]interface{}
	expirationTime map[string]time.Time
	mutex          sync.RWMutex
	cleanupTick    time.Duration
}

// Initialize the cache
func NewCache(cleanupTick time.Duration) *Cache {
	cache := &Cache{
		data:           make(map[string]interface{}),
		expirationTime: make(map[string]time.Time),
		cleanupTick:    cleanupTick,
	}
	go cache.startCleanup()
	return cache
}

// Start the cleanup routine
func (cache *Cache) startCleanup() {
	ticker := time.NewTicker(cache.cleanupTick)
	for range ticker.C {
		cache.cleanupCache()
	}
}

// Cleanup the cache
func (cache *Cache) cleanupCache() {
	currentTime := time.Now()
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, expirationTime := range cache.expirationTime {
		if currentTime.After(expirationTime) {
			delete(cache.data, key)
			delete(cache.expirationTime, key)
		}
	}
}

// Set a value in the cache
func (cache *Cache) set(key string, value interface{}, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data[key] = value
	cache.expirationTime[key] = time.Now().Add(ttl)
}

// Get a value from the cache
func (cache *Cache) get(key string) (interface{}, bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	value, exists := cache.data[key]
	return value, exists
}

func main() {
	cache := NewCache(1 * time.Second)

	cache.set("1", "Apple", 3*time.Second)
	cache.set("2", "Kiwi", 4*time.Second)
	cache.set("1", "Strawberry", 6*time.Second)

	for key, _ := range cache.data {
		value, exists := cache.get(key)
		if exists {
			fmt.Println("value: ", value)
		}
	}

	time.Sleep(7 * time.Second)
	value, exists := cache.get("1")

	if exists {
		fmt.Println(value)
	} else {
		fmt.Println("Not found")
	}
}
