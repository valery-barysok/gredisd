package model

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestIntegrationForAllCommandsAtOnce(t *testing.T) {
	RegisterTestingT(t)

	dbModel := newDBModel(0)
	if Expect(dbModel).ToNot(Equal(nil)) {
		//pingRes, err := dbModel.Ping()
		//Expect(err).ToNot(HaveOccurred())
		//Expect(pingRes).To(Equal("PONG"))
		//
		//msg := "test echo cmd"
		//bulkRes, err := dbModel.PingMsg(msg)
		//Expect(err).ToNot(HaveOccurred())
		//Expect(bulkRes).To(BeEquivalentTo(msg))
		//
		//bulkRes, err = dbModel.Echo(msg)
		//Expect(err).ToNot(HaveOccurred())
		//Expect(bulkRes).To(BeEquivalentTo(msg))
		//
		//list, err := dbModel.Command()
		//Expect(err).ToNot(HaveOccurred())
		//Expect(list).ToNot(Equal(nil))

		// Valid regexp
		list, err := dbModel.Keys([]byte(".*"))
		Expect(err).ToNot(HaveOccurred())
		if Expect(list).ToNot(Equal(nil)) {
			Expect(len(list)).To(Equal(0))
		}

		// Invalid regexp
		list, err = dbModel.Keys([]byte(")"))
		Expect(err).To(HaveOccurred())
		if Expect(list).ToNot(Equal(nil)) {
			Expect(len(list)).To(Equal(0))
		}

		key := []byte("key")
		exists := dbModel.Exists(key)
		Expect(exists).To(Equal(0))

		keyValue := []byte("key_value")
		dbModel.Set(key, keyValue)

		exists = dbModel.Exists(key)
		Expect(exists).To(Equal(1))

		exists = dbModel.Exists(key, key, key)
		Expect(exists).To(Equal(3))

		value, err := dbModel.Get(key)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(keyValue))

		keys, err := dbModel.Keys([]byte(".*"))
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(key))
		}

		cnt := dbModel.Del(key)
		Expect(cnt).To(Equal(1))

		exists = dbModel.Exists(key)
		Expect(exists).To(Equal(0))

		listKey := []byte("list_key")
		listKeyValue1 := []byte("list_key_value1")
		listKeyValue2 := []byte("list_key_value2")
		listKeyValue3 := []byte("list_key_value3")
		listKeyValue4 := []byte("list_key_value4")
		listKeyValue5 := []byte("list_key_value5")
		listKeyValue6 := []byte("list_key_value6")
		listKeyValue7 := []byte("list_key_value7")
		listKeyValue8 := []byte("list_key_value8")

		cnt, err = dbModel.LPush(listKey, listKeyValue1, listKeyValue2, listKeyValue3)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(3))

		cnt, err = dbModel.RPush(listKey, listKeyValue4, listKeyValue5, listKeyValue6)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(6))

		value, err = dbModel.LPop(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue3))

		value, err = dbModel.RPop(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue6))

		l, err := dbModel.LLen(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(l).To(Equal(4))

		cnt, err = dbModel.LInsert(listKey, insertBefore, listKeyValue2, listKeyValue7)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(5))

		cnt, err = dbModel.LInsert(listKey, insertAfter, listKeyValue2, listKeyValue8)
		Expect(err).ToNot(HaveOccurred())
		Expect(cnt).To(Equal(6))

		value, err = dbModel.LIndex(listKey, []byte("0"))
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue7))

		value, err = dbModel.LIndex(listKey, []byte("2"))
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(listKeyValue8))

		values, err := dbModel.LRange(listKey, []byte("0"), []byte("-1"))
		Expect(err).ToNot(HaveOccurred())
		if Expect(values).To(HaveLen(6)) {
			Expect(values).To(Equal([]interface{}{
				listKeyValue7,
				listKeyValue2,
				listKeyValue8,
				listKeyValue1,
				listKeyValue4,
				listKeyValue5,
			}))
		}

		keys, err = dbModel.Keys([]byte(".*"))
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(listKey))
		}

		cnt = dbModel.Del(listKey)
		Expect(cnt).To(Equal(1))

		exists = dbModel.Exists(listKey)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		dictKey := []byte("dict_key")
		dictKeyField := []byte("dict_key_field")
		dictKeyField2 := []byte("dict_key_field2")
		dictKeyFieldValue := []byte("dict_key_field_value")
		dictKeyFieldValue2 := []byte("dict_key_field_value2")

		exists, err = dbModel.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		inserted, err := dbModel.HSet(dictKey, dictKeyField, dictKeyFieldValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(inserted).To(Equal(1))

		updated, err := dbModel.HSet(dictKey, dictKeyField, dictKeyFieldValue)
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).To(Equal(0))

		inserted, err = dbModel.HSet(dictKey, dictKeyField2, dictKeyFieldValue2)
		Expect(err).ToNot(HaveOccurred())
		Expect(inserted).To(Equal(1))

		updated, err = dbModel.HSet(dictKey, dictKeyField2, dictKeyFieldValue2)
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).To(Equal(0))

		exists, err = dbModel.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(1))

		value, err = dbModel.HGet(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(value).To(BeEquivalentTo(dictKeyFieldValue))

		deleted, err := dbModel.HDel(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(deleted).To(Equal(1))

		exists, err = dbModel.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))

		keys, err = dbModel.Keys([]byte(".*"))
		Expect(err).ToNot(HaveOccurred())
		if Expect(keys).To(HaveLen(1)) {
			Expect(keys[0]).To(BeEquivalentTo(dictKey))
		}

		cnt = dbModel.Del(dictKey)
		Expect(cnt).To(Equal(1))

		exists = dbModel.Exists(dictKey)
		Expect(exists).To(Equal(0))

		exists, err = dbModel.HExists(dictKey, dictKeyField)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(Equal(0))
	}
}
