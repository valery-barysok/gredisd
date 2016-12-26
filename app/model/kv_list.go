package model

import (
	"bytes"
	"container/list"
)

type lrPush func(list *list.List, v interface{}) *list.Element
type lrPop func(list *list.List) []byte

func newKeyValueList() *keyValue {
	return &keyValue{
		kvType: kvListType,
		list:   list.New(),
	}
}

func (kv *kvModel) LPush(key []byte, values ...[]byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.lpush(key, values...)
}

func (kv *kvModel) RPush(key []byte, values ...[]byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.rpush(key, values...)
}

func (kv *kvModel) LPop(key []byte) ([]byte, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.lpop(key)
}

func (kv *kvModel) RPop(key []byte) ([]byte, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.rpop(key)
}

func (kv *kvModel) LLen(key []byte) (int, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.llen(key)
}

func (kv *kvModel) LInsert(key []byte, before bool, pivot []byte, value []byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.linsert(key, before, pivot, value)
}

func (kv *kvModel) LIndex(key []byte, index int) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.lindex(key, index)
}

func (kv *kvModel) LRange(key []byte, start int, stop int) ([]interface{}, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.lrange(key, start, stop)
}

func (kv *kvModel) lrpush(push lrPush, key []byte, values ...[]byte) (int, error) {
	k := string(key)
	val, exists := kv.tryGet(k)
	if exists {
		if val.kvType != kvListType {
			return 0, errWrongType
		}
	} else {
		val = newKeyValueList()
		kv.storage[k] = val
	}

	for _, value := range values {
		push(val.list, value)
	}

	return val.list.Len(), nil
}

func (kv *kvModel) lpush(key []byte, values ...[]byte) (int, error) {
	return kv.lrpush(func(list *list.List, v interface{}) *list.Element {
		return list.PushFront(v)
	}, key, values...)
}

func (kv *kvModel) rpush(key []byte, values ...[]byte) (int, error) {
	return kv.lrpush(func(list *list.List, v interface{}) *list.Element {
		return list.PushBack(v)
	}, key, values...)
}

func (kv *kvModel) lrpop(pop lrPop, key []byte) ([]byte, error) {
	k := string(key)
	val, exists := kv.tryGet(k)
	if exists {
		if val.kvType != kvListType {
			return nil, errWrongType
		}
	} else {
		return nil, nil
	}

	e := pop(val.list)
	if val.list.Len() == 0 {
		delete(kv.storage, k)
	}

	return e, nil
}

func (kv *kvModel) lpop(key []byte) ([]byte, error) {
	return kv.lrpop(func(list *list.List) []byte {
		return list.Remove(list.Front()).([]byte)
	}, key)
}

func (kv *kvModel) rpop(key []byte) ([]byte, error) {
	return kv.lrpop(func(list *list.List) []byte {
		return list.Remove(list.Back()).([]byte)
	}, key)
}

func (kv *kvModel) llen(key []byte) (int, error) {
	val, exists := kv.tryGet(string(key))
	if exists {
		if val.kvType != kvListType {
			return 0, errWrongType
		}

		return val.list.Len(), nil
	}

	return 0, nil
}

func (kv *kvModel) linsert(key []byte, before bool, pivot []byte, value []byte) (int, error) {
	val, exists := kv.tryGet(string(key))
	if exists {
		if val.kvType != kvListType {
			return 0, errWrongType
		}

		for it := val.list.Front(); it != nil; it = it.Next() {
			if bytes.Equal(it.Value.([]byte), pivot) {
				if before {
					val.list.InsertBefore(value, it)
				} else {
					val.list.InsertAfter(value, it)
				}

				return val.list.Len(), nil
			}
		}

		return -1, nil
	}

	return 0, nil
}

func (kv *kvModel) lindex(key []byte, index int) ([]byte, error) {
	val, exists := kv.tryGet(string(key))
	if exists {
		if val.kvType != kvListType {
			return nil, errWrongType
		}

		l := val.list.Len()
		index = leftIndex(index, l)
		if 0 <= index && index < l {
			for it := val.list.Front(); it != nil; it = it.Next() {
				if index == 0 {
					return it.Value.([]byte), nil
				}
				index--
			}
		}
	}

	return nil, nil
}

func (kv *kvModel) lrange(key []byte, start int, stop int) ([]interface{}, error) {
	val, exists := kv.tryGet(string(key))
	if exists {
		if val.kvType != kvListType {
			return nil, errWrongType
		}

		l := val.list.Len()
		left := leftIndex(start, l)
		right := rightIndex(stop, l)
		if left < right {
			values := make([]interface{}, 0, right-left)
			for i, it := 0, val.list.Front(); it != nil; it = it.Next() {
				if i >= right {
					break
				}
				if i >= left {
					values = append(values, it.Value)
				}
				i++
			}
			return values, nil
		}
	}

	return make([]interface{}, 0), nil
}

func leftIndex(left int, total int) int {
	if left < 0 {
		return max(total+left, 0)
	}
	return left
}

func rightIndex(right int, total int) int {
	if right < 0 {
		return min(total+right+1, total)
	}
	return right + 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
