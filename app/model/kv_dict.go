package model

func newKeyValueDict() *keyValue {
	return &keyValue{
		kvType: kvDictType,
		dict:   make(map[string]string),
	}
}

func (kv *kvModel) HSet(key []byte, field []byte, value []byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.hset(key, field, value)
}

func (kv *kvModel) HGet(key []byte, field []byte) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.hget(key, field)
}

func (kv *kvModel) HDel(key []byte, fields ...[]byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.hdel(key, fields...)
}

func (kv *kvModel) HLen(key []byte, fields ...[]byte) (int, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.hlen(key)
}

func (kv *kvModel) HExists(key []byte, field []byte) (int, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.hexists(key, field)
}

func (kv *kvModel) hset(key []byte, field []byte, value []byte) (int, error) {
	k := string(key)
	val, exists := kv.storage[k]
	if exists {
		if val.kvType != kvDictType {
			return 0, errWrongType
		}
	} else {
		val = newKeyValueDict()
		kv.storage[k] = val
	}

	f := string(field)
	_, ok := val.dict[f]
	val.dict[f] = string(value)
	if ok {
		return 0, nil
	}
	return 1, nil
}

func (kv *kvModel) hget(key []byte, field []byte) ([]byte, error) {
	k := string(key)
	val, exists := kv.storage[k]
	if !exists {
		return nil, nil
	}

	if val.kvType != kvDictType {
		return nil, errWrongType
	}

	v, exists := val.dict[string(field)]
	if exists {
		return []byte(v), nil
	}
	return nil, nil
}

func (kv *kvModel) hdel(key []byte, fields ...[]byte) (int, error) {
	k := string(key)
	val, exists := kv.storage[k]
	if !exists {
		return 0, nil
	}

	if val.kvType != kvDictType {
		return 0, errWrongType
	}

	cnt := 0
	for _, field := range fields {
		f := string(field)
		_, exists = val.dict[f]
		if exists {
			cnt++
		}
		delete(val.dict, f)
	}

	return cnt, nil
}

func (kv *kvModel) hlen(key []byte) (int, error) {
	k := string(key)
	val, exists := kv.storage[k]
	if exists {
		if val.kvType != kvDictType {
			return 0, errWrongType
		}

		return val.list.Len(), nil
	}

	return 0, nil
}

func (kv *kvModel) hexists(key []byte, field []byte) (int, error) {
	k := string(key)
	val, exists := kv.storage[k]
	if exists {
		if val.kvType != kvDictType {
			return 0, errWrongType
		}

		_, exists := val.dict[string(field)]
		if exists {
			return 1, nil
		}
	}

	return 0, nil
}
