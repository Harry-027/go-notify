package handler

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type CacheConn struct {
	pool *redis.Pool
}

var Cache CacheConn

func RedisPoolInit() {
	Cache = CacheConn{
		pool: &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", "localhost:6379")
			},
		},
	}
}

func (cache *CacheConn) setDetails(key, value string) error {
	conn := cache.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("SET", key, value)
	if err != nil {
		log.Println("Error while setting key: ", err)
		return err
	}
	log.Println("Cache server reply on key set: ", reply)

	reply, err = conn.Do("EXPIRE", key, 3600)
	if err != nil {
		log.Println("Error while setting key expiry: ", err)
		return err
	}
	log.Println("Cache server reply on key expire set: ", reply)
	return nil
}

func (cache *CacheConn) getDetails(key string) (string, error) {
	conn := cache.pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Println("An error occurred while fetching key from cache", err.Error())
		return "", err
	}
	return reply, nil
}

func (cache *CacheConn) ifExistsInCache(key string) (bool, error) {
	conn := cache.pool.Get()
	defer conn.Close()
	exists, err := redis.Int(conn.Do("EXISTS", key))
	log.Println("Exists in cache: ", exists)
	if err != nil {
		log.Println("An error occurred while checking if the key exists in cache", err.Error())
		return false, err
	}
	if exists == 0 {
		return false, errors.New("key doesn't exists")
	}
	return true, nil
}

func (cache *CacheConn) deleteKey(key string) (bool, error) {
	conn := cache.pool.Get()
	defer conn.Close()
	_, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		log.Println("An error occurred while deleting key from cache: ", err.Error())
		return false, err
	}
	return true, nil

}

func (cache *CacheConn) setExpiry(key string, sec int) {
	conn := cache.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("EXPIRE", key, sec)
	if err != nil {
		log.Println("Error while setting key expiry: ", err)
	}
	log.Println("Cache server reply on key expire set: ", reply)
}
