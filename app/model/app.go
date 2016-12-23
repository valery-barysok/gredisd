package model

import (
	"container/list"
	"errors"
	"strconv"
	"sync"
)

var errInvalidDBIndex = errors.New("ERR invalid DB index")

type AppModel struct {
	mu        sync.Mutex
	databases int
	dbs       map[int]*DBModel
	commands  *list.List
}

func NewAppModel(databases int) *AppModel {
	return &AppModel{
		databases: databases,
		dbs:       make(map[int]*DBModel),
		commands:  list.New(),
	}
}

func (model *AppModel) Select(index string) (*DBModel, error) {
	ind, err := strconv.Atoi(index)
	if err != nil {
		return nil, errInvalidDBIndex
	}

	return model.SelectIndex(ind)
}

func (model *AppModel) SelectIndex(index int) (*DBModel, error) {
	if 0 > index || index >= model.databases {
		return nil, errInvalidDBIndex
	}

	db, ok := model.dbs[index]
	if ok {
		return db, nil
	}

	model.mu.Lock()
	defer model.mu.Unlock()

	db, ok = model.dbs[index]
	if ok {
		return db, nil
	}

	db = newDBModel(index)
	model.dbs[db.index] = db

	return db, nil
}

func (model *AppModel) AddCmd(cmd string) {
	model.commands.PushBack([]byte(cmd))
}

func (model *AppModel) Commands() []interface{} {
	cmds := make([]interface{}, 0, model.commands.Len())
	for it := model.commands.Front(); it != nil; it = it.Next() {
		cmds = append(cmds, it.Value)
	}
	return cmds
}
