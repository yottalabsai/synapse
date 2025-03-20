package utils

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

type Lock struct {
	Key     string
	Token   string
	Conn    redis.Conn
	Timeout int
}

func (lock *Lock) TryLock() (ok bool, err error) {
	_, err = redis.String(lock.Conn.Do("SET", lock.Key, lock.Token, "EX", lock.Timeout, "NX"))
	if errors.Is(err, redis.ErrNil) {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *Lock) Unlock() (err error) {
	_, err = lock.Conn.Do("del", lock.Key)
	return
}
