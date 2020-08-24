//Package cache the cache interface
package cache

import (
	"time"
)

//Cache the Cache interface defination
type Cache interface {
	PaddingKey(key string) string

	Set(key string, val interface{}, duration time.Duration) error
	SetString(key, val string, duration time.Duration) error
	SetInt(key string, val int64, duration time.Duration) error
	SetObject(key string, val interface{}, duration time.Duration) error

	Get(key string) (interface{}, error)
	IsSet(key string) (bool, error)
	Remove(key string) (bool, error)
	GetString(key string) (string, error)
	GetInt(key string) (int64, error)
	GetObject(key string, val interface{}) (bool, error)

	GetMap(key string) (map[string]interface{}, error)

	SafetyIncr(key string, step int64) (bool, error)
}

//SyncLocker a sync locker interface
type SyncLocker interface {
	GetLock(key string, expire time.Duration) bool
	ReleaseLock(key string) error
}
