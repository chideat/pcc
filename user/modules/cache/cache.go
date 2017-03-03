package cache

import (
	"errors"
	"fmt"

	"github.com/chideat/glog"
	"github.com/garyburd/redigo/redis"
	"stathat.com/c/consistent"
)

// connection cache
// not goroutine safe
type Cache struct {
	conns      map[string]*redis.Pool
	consistent *consistent.Consistent
}

func (c *Cache) Get(key string) (redis.Conn, error) {
	if name, err := c.consistent.Get(key); err == nil {
		return c.conns[name].Get(), nil
	} else {
		glog.Error(err)
		return nil, err
	}
}

func (c *Cache) GetByName(name string) (redis.Conn, error) {
	if c, ok := c.conns[name]; ok {
		return c.Get(), nil
	} else {
		return nil, errors.New(fmt.Sprintf("INVALID name %s.", name))
	}
}

func (c *Cache) Do(cmd string, key string, args ...interface{}) (interface{}, error) {
	name, err := c.consistent.Get(key)
	if err != nil {
		return nil, err
	}

	conn := c.conns[name].Get()
	reply, err := conn.Do(cmd, append([]interface{}{key}, args...)...)
	conn.Close()

	return reply, err
}

func NewCache(addrs map[string]string) *Cache {
	cache := new(Cache)
	cache.conns = map[string]*redis.Pool{}
	cache.consistent = consistent.New()

	for key, addr := range addrs {
		cache.conns[key] = func(addr string) *redis.Pool {
			return redis.NewPool(func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}

				if _, err = conn.Do("PING"); err != nil {
					conn.Close()
					return nil, err
				}
				return conn, err
			}, 1024)
		}(addr)
		cache.consistent.Add(key)
	}

	return cache
}
