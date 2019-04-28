package main

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

//lru-k   k=2

type CacheNode struct {
	Key, Value interface{}
	Hit        int
}

func (cnode *CacheNode) NewCacheNode(k, v interface{}) *CacheNode {
	return &CacheNode{k, v, 0}
}

type LRUK2 struct {
	L1 *LRUCache
	L2 *LRUCache
}

type LRUCache struct {
	Capacity int
	HitCount int
	dlist    *list.List
	cacheMap map[interface{}]*list.Element
	next     *LRUCache
}

func NewLRUK2(cap1, cap2, hitcount int) *LRUK2 {

	l2 := &LRUCache{
		Capacity: cap2,
		HitCount: hitcount,
		dlist:    list.New(),
		next:     nil,
		cacheMap: make(map[interface{}]*list.Element)}

	l1 := &LRUCache{
		Capacity: cap1,
		HitCount: hitcount,
		dlist:    list.New(),
		next:     l2,
		cacheMap: make(map[interface{}]*list.Element)}

	return &LRUK2{
		L1: l1,
		L2: l2,
	}

}

func NewLRUCache(cap, hit int) *LRUCache {
	return &LRUCache{
		Capacity: cap,
		HitCount: hit,
		dlist:    list.New(),
		cacheMap: make(map[interface{}]*list.Element)}
}

func (lru *LRUCache) Size() int {
	return lru.dlist.Len()
}

func (lru *LRUCache) Set(k, v interface{}) error {

	if lru.dlist == nil {
		return errors.New("LRUCache结构体未初始化.")
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(pElement)
		pElement.Value.(*CacheNode).Value = v
		pElement.Value.(*CacheNode).Hit++
		if pElement.Value.(*CacheNode).Hit >= lru.HitCount && lru.next != nil {
			lru.next.Set(k, v)
			lru.Remove(k)
		}
		return nil
	}

	newElement := lru.dlist.PushFront(&CacheNode{k, v, 1})
	lru.cacheMap[k] = newElement

	if lru.dlist.Len() > lru.Capacity {
		//移掉最后一个
		lastElement := lru.dlist.Back()
		if lastElement == nil {
			return nil
		}
		cacheNode := lastElement.Value.(*CacheNode)
		delete(lru.cacheMap, cacheNode.Key)
		lru.dlist.Remove(lastElement)
	}
	return nil
}

func (lru *LRUCache) Get(k interface{}) (v interface{}, ret bool, err error) {

	if lru.cacheMap == nil {
		return v, false, errors.New("LRUCache结构体未初始化.")
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(pElement)
		return pElement.Value.(*CacheNode).Value, true, nil
	}
	return v, false, nil
}

func (lru *LRUCache) Remove(k interface{}) bool {

	if lru.cacheMap == nil {
		return false
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		cacheNode := pElement.Value.(*CacheNode)
		delete(lru.cacheMap, cacheNode.Key)
		lru.dlist.Remove(pElement)
		return true
	}
	return false
}

func (lru *LRUCache) PrintList(w io.Writer) {

	if lru.dlist.Len() != 0 {
		f := lru.dlist.Front()
		for f != lru.dlist.Back() {
			fmt.Fprintln(w, f.Value.(*CacheNode).Value, f.Value.(*CacheNode).Hit)
			f = f.Next()
		}
	}

}

func main() {

	data, _ := ioutil.ReadFile("data.txt")
	number := bytes.Split(data, []byte(" "))

	lru2 := NewLRUK2(len(data)*20/100, len(data)*80/100, 3)

	for index, _ := range number {
		intdata, _ := strconv.Atoi(string(number[index]))

		_, l2get, _ := lru2.L2.Get(intdata)

		if l2get {
			lru2.L2.Set(intdata, intdata)
		} else {
			lru2.L1.Set(intdata, intdata)
		}

	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("l2:\n")
	lru2.L2.PrintList(buf)
	buf.WriteString("l1:\n")

	lru2.L1.PrintList(buf)
	ioutil.WriteFile("lru.out", buf.Bytes(), 0755)

}
