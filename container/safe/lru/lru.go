package lru

import (
	"container/list"
	"sync"
)

type node struct {
	k interface{}
	v interface{}
}

type LRU struct {
	mutex   sync.RWMutex
	maxLen  int // if maxLen<=0 then no limit
	lruList *list.List
	datas   map[interface{}]*list.Element
}

func New(maxLen int) *LRU {
	return &LRU{
		maxLen:  maxLen,
		lruList: list.New(),
		datas:   make(map[interface{}]*list.Element),
	}
}

func (lc *LRU) PushFront(k, v interface{}) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if e, ok := lc.datas[k]; ok {
		e.Value.(*node).v = v
		lc.lruList.MoveToFront(e)
	} else {
		e := lc.lruList.PushFront(&node{k, v})
		lc.datas[k] = e
	}
	if lc.maxLen > 0 && lc.lruList.Len() > lc.maxLen {
		lc.popTail()
	}
}

func (lc *LRU) PopTail() (interface{}, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.popTail()
}

func (lc *LRU) Get(k interface{}) (interface{}, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if e, ok := lc.datas[k]; ok {
		lc.lruList.MoveToFront(e)
		return e.Value.(*node).v, ok
	}
	return nil, false
}

func (lc *LRU) Del(k interface{}) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.del(k)
}

func (lc *LRU) Len() int {
	lc.mutex.RLock()
	defer lc.mutex.RUnlock()
	return lc.lruList.Len()
}

func (lc *LRU) popTail() (interface{}, bool) {
	e := lc.lruList.Back()
	if e != nil {
		lc.del(e.Value.(*node).k)
		return e.Value.(*node).v, true
	}
	return nil, false
}

func (lc *LRU) del(k interface{}) {
	if e, ok := lc.datas[k]; ok {
		delete(lc.datas, k)
		lc.lruList.Remove(e)
	}
}
