package model

import (
	"container/list"
	"errors"
	"regexp"
	"sync"
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

func (kv *kvModel) Expire(key []byte, value []byte) {
	//kv.mu.Lock()
	//defer kv.mu.Unlock()
	//kv.set(key, value)
}

func (kv *kvModel) keys(pattern []byte) ([]interface{}, error) {
	re, err := regexp.CompilePOSIX(string(pattern))
	if err != nil {
		return nil, err
	}

	lst := list.New()
	for k := range kv.storage {
		if re.MatchString(k) {
			lst.PushBack([]byte(k))
		}
	}

	keys := make([]interface{}, 0, lst.Len())
	for e := lst.Front(); e != nil; e = lst.Front() {
		keys = append(keys, lst.Remove(e))
	}
	return keys, nil
}

func (kv *kvModel) keyExists(key []byte) bool {
	_, exists := kv.storage[string(key)]
	return exists
}

func (kv *kvModel) exists(keys ...[]byte) int {
	cnt := 0
	for _, key := range keys {
		if kv.keyExists(key) {
			cnt++
		}
	}
	return cnt
}
