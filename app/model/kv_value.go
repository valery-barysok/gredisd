package model

func newKeyValue(s []byte) *keyValue {
	return &keyValue{
		kvType: kvType,
		value:  s,
	}
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

func (kv *KVModel) set(key []byte, value []byte) {
	kv.storage[string(key)] = newKeyValue(value)
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
