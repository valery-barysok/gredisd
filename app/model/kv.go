package model

import (
	"container/list"
	"errors"
	"regexp"
	"sync"
)

var errWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

type valueItem struct {
	kvType byte
	value  []byte
	list   *list.List
	dict   map[string]string
	ttl    int64
}

func newKVType(s []byte) *valueItem {
	return &valueItem{
		kvType: kvType,
		value:  s,
	}
}

type KVModel struct {
	mu      sync.RWMutex
	storage map[string]*valueItem
	db      *DBModel
}

func newKVModel(db *DBModel) *KVModel {
	return &KVModel{
		storage: make(map[string]*valueItem),
		db:      db,
	}
}

func (kv *KVModel) Keys(pattern []byte) ([]interface{}, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.keys(pattern)
}

func (kv *KVModel) Set(key []byte, value []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.set(key, value)
}

func (kv *KVModel) Get(key []byte) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.get(key)
}

func (kv *KVModel) Del(keys ...[]byte) int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.delKeys(keys...)
}

func (kv *KVModel) Exists(keys ...[]byte) int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.exists(keys...)
}

func (kv *KVModel) Expire(key []byte, value []byte) {
	//kv.mu.Lock()
	//defer kv.mu.Unlock()
	//kv.set(key, value)
}

func (kv *KVModel) keys(pattern []byte) ([]interface{}, error) {
	re, err := regexp.CompilePOSIX(string(pattern))
	if err != nil {
		return nil, err
	}

	list := list.New()
	for k := range kv.storage {
		if re.MatchString(k) {
			list.PushBack([]byte(k))
		}
	}

	keys := make([]interface{}, 0, list.Len())
	for e := list.Front(); e != nil; e = list.Front() {
		keys = append(keys, list.Remove(e))
	}
	return keys, nil
}

func (kv *KVModel) set(key []byte, value []byte) {
	kv.storage[string(key)] = newKVType(value)
}

func (kv *KVModel) get(key []byte) ([]byte, error) {
	val, exists := kv.storage[string(key)]
	if exists {
		if val.kvType != kvType {
			return nil, errWrongType
		}
		return val.value, nil
	}

	return nil, nil
}

func (kv *KVModel) delKeys(keys ...[]byte) int {
	cnt := 0
	for _, key := range keys {
		cnt += kv.del(key)
	}
	return cnt
}

func (kv *KVModel) del(key []byte) int {
	if kv.keyExists(key) {
		delete(kv.storage, string(key))
		return 1
	}
	return 0
}

func (kv *KVModel) keyExists(key []byte) bool {
	_, exists := kv.storage[string(key)]
	return exists
}

func (kv *KVModel) exists(keys ...[]byte) int {
	cnt := 0
	for _, key := range keys {
		if kv.keyExists(key) {
			cnt++
		}
	}
	return cnt
}
