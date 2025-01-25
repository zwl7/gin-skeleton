package utils

import (
	"net/url"
	"time"
)

type CacheInterface interface {
	Get(key string) interface{}
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}

type ApiCacheInterface interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
}

// ParamsInterface ...
type ParamsInterface interface {
	// 返回参数列表
	Params() (url.Values, error)
}

type LogerInterface interface {
	//SetLevel(level uint32)
	//GetLevel() uint32
	Trace(args ...interface{})
	Tracef(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
}

type RedisInterface interface {
	Lock(key string, timeout ...time.Duration) bool
	UnLock(key string) int64
}
