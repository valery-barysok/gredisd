package model

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	errSyntax         = errors.New("syntax error")
	errInvalidInteger = errors.New("value is not an integer or out of range")
)

var (
	insertBefore = []byte("BEFORE")
	insertAfter  = []byte("AFTER")
)

type DBModel struct {
	index int
	kv    *kvModel
}

func newDBModel(index int) *DBModel {
	model := &DBModel{
		index: index,
	}
	model.kv = newKVModel()
	return model
}

func (db *DBModel) Keys(pattern []byte) ([]interface{}, error) {
	return db.kv.Keys(pattern)
}

func (db *DBModel) Set(key []byte, value []byte) {
	db.kv.Set(key, value)
}

func (db *DBModel) Get(key []byte) ([]byte, error) {
	return db.kv.Get(key)
}

func (db *DBModel) Del(keys ...[]byte) int {
	return db.kv.Del(keys...)
}

func (db *DBModel) Exists(keys ...[]byte) int {
	return db.kv.Exists(keys...)
}

func (db *DBModel) Expire(key []byte, seconds []byte) int {
	// TODO: implement Expire
	return 0
}

func (db *DBModel) LPush(key []byte, values ...[]byte) (int, error) {
	return db.kv.LPush(key, values...)
}

func (db *DBModel) RPush(key []byte, values ...[]byte) (int, error) {
	return db.kv.RPush(key, values...)
}

func (db *DBModel) LPop(key []byte) ([]byte, error) {
	return db.kv.LPop(key)
}

func (db *DBModel) RPop(key []byte) ([]byte, error) {
	return db.kv.RPop(key)
}

func (db *DBModel) LLen(key []byte) (int, error) {
	return db.kv.LLen(key)
}

func (db *DBModel) LInsert(key []byte, place []byte, pivot []byte, value []byte) (int, error) {
	before := bytes.EqualFold(place, insertBefore)
	after := bytes.EqualFold(place, insertAfter)
	if !before && !after {
		return -1, errSyntax
	}

	return db.LInsertN(key, before, pivot, value)
}

func (db *DBModel) LInsertN(key []byte, before bool, pivot []byte, value []byte) (int, error) {
	return db.kv.LInsert(key, before, pivot, value)
}

func (db *DBModel) LIndex(key []byte, index []byte) ([]byte, error) {
	ind, err := strconv.Atoi(string(index))
	if err != nil {
		return nil, errInvalidInteger
	}

	return db.LIndexN(key, ind)
}

func (db *DBModel) LIndexN(key []byte, index int) ([]byte, error) {
	return db.kv.LIndex(key, index)
}

func (db *DBModel) LRange(key []byte, start []byte, stop []byte) ([]interface{}, error) {
	s, err := strconv.Atoi(string(start))
	if err != nil {
		return nil, errInvalidInteger
	}

	e, err := strconv.Atoi(string(stop))
	if err != nil {
		return nil, errInvalidInteger
	}

	return db.LRangeN(key, s, e)
}

func (db *DBModel) LRangeN(key []byte, start int, stop int) ([]interface{}, error) {
	return db.kv.LRange(key, start, stop)
}

func (db *DBModel) HSet(key []byte, field []byte, value []byte) (int, error) {
	return db.kv.HSet(key, field, value)
}

func (db *DBModel) HGet(key []byte, field []byte) ([]byte, error) {
	return db.kv.HGet(key, field)
}

func (db *DBModel) HDel(key []byte, fields ...[]byte) (int, error) {
	return db.kv.HDel(key, fields...)
}

func (db *DBModel) HLen(key []byte) (int, error) {
	return db.kv.HLen(key)
}

func (db *DBModel) HExists(key []byte, field []byte) (int, error) {
	return db.kv.HExists(key, field)
}
