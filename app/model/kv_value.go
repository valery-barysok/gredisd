package model

func newKeyValue(s []byte) *keyValue {
	return &keyValue{
		kvType: kvType,
		value:  s,
	}
}

func (kv *kvModel) Set(key []byte, value []byte) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.set(key, value)
}

func (kv *kvModel) Get(key []byte) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.get(key)
}

func (kv *kvModel) Del(keys ...[]byte) int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	return kv.delKeys(keys...)
}

func (kv *kvModel) set(key []byte, value []byte) {
	kv.storage[string(key)] = newKeyValue(value)
}

func (kv *kvModel) get(key []byte) ([]byte, error) {
	val, exists := kv.tryGet(string(key))
	if exists {
		if val.kvType != kvType {
			return nil, errWrongType
		}
		return val.value, nil
	}

	return nil, nil
}

func (kv *kvModel) delKeys(keys ...[]byte) int {
	cnt := 0
	for _, key := range keys {
		cnt += kv.del(key)
	}
	return cnt
}

func (kv *kvModel) del(key []byte) int {
	if kv.keyExists(key) {
		delete(kv.storage, string(key))
		return 1
	}
	return 0
}
