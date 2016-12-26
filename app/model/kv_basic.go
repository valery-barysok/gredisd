package model

import (
	"container/list"
	"errors"
	"regexp"
	"sync"
	"time"
)

const (
	kvType     byte = 1
	kvListType byte = 2
	kvDictType byte = 3
)

var errWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

type keyValue struct {
	kvType byte
	value  []byte
	list   *list.List
	dict   map[string]string
	ttl    int64
}

type kvModel struct {
	mu      sync.RWMutex
	storage map[string]*keyValue
}

func newKVModel() *kvModel {
	return &kvModel{
		storage: make(map[string]*keyValue),
	}
}

func (kv *kvModel) Keys(pattern []byte) ([]interface{}, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.keys(pattern)
}

func (kv *kvModel) Exists(keys ...[]byte) int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.exists(keys...)
}

func (kv *kvModel) Expire(key []byte, ttl int64) int {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.expire(key, ttl)
}

func (kv *kvModel) keys(pattern []byte) ([]interface{}, error) {
	re, err := regexp.CompilePOSIX(string(pattern))
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	lst := list.New()
	for k := range kv.storage {
		if re.MatchString(k) {
			_, exists := kv.tryGetN(k, now)
			if exists {
				lst.PushBack([]byte(k))
			}
		}
	}

	keys := make([]interface{}, 0, lst.Len())
	for e := lst.Front(); e != nil; e = lst.Front() {
		keys = append(keys, lst.Remove(e))
	}
	return keys, nil
}

func (kv *kvModel) keyExists(key []byte) bool {
	return kv.keyExistsN(key, time.Now().Unix())
}

func (kv *kvModel) keyExistsN(key []byte, now int64) bool {
	_, exists := kv.tryGetN(string(key), now)
	return exists
}

func (kv *kvModel) exists(keys ...[]byte) int {
	now := time.Now().Unix()

	cnt := 0
	for _, key := range keys {
		if kv.keyExistsN(key, now) {
			cnt++
		}
	}
	return cnt
}

func (kv *kvModel) expire(key []byte, ttl int64) int {
	now := time.Now().Unix()

	val, exists := kv.tryGetN(string(key), now)
	if !exists {
		return 0
	}

	val.ttl = ttl + now
	return 1
}

func (kv *kvModel) tryGet(key string) (*keyValue, bool) {
	return kv.tryGetN(key, time.Now().Unix())
}

func (kv *kvModel) tryGetN(key string, now int64) (*keyValue, bool) {
	val, exists := kv.storage[key]
	if !exists {
		return nil, false
	}

	if isExpired(val, now) {
		delete(kv.storage, key)
		return nil, false
	}

	return val, true
}

func isExpired(val *keyValue, now int64) bool {
	return val.ttl != 0 && val.ttl-now <= 0
}
