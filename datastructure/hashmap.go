package datastruct

//a fix-size memory hashmap,use lru
//use slice simulate hashtable

import (
	"fmt"
	"github.com/pandaychen/goes-wrapper/hashalgo"
	"strconv"
	"sync"
	"time"
)

const (
	DEFAULT_BUCKET_LEN = 8
	DEFAULT_BUCKET_NUM = 1024
)

//unit-node
type SHashNode struct {
	key    string
	value  interface{}
	vistor int64
}

//Manager
type SHashTable struct {
	nodes          []*SHashNode
	collisionNodes []*SHashNode
	bucketLen      int //桶深
	bucketNum      int //桶个数
	sync.RWMutex
}

//init
func (h *SHashTable) Init(blen, bnum int) *SHashTable {
	if h == nil {
		h = new(SHashTable)
	}

	if blen == 0 || bnum == 0 {
		h.bucketLen = DEFAULT_BUCKET_LEN
		h.bucketNum = int(float64(h.bucketLen)*1.2)/h.bucketLen + 1 //suggest to be a prime
		h.nodes = make([]*SHashNode, h.bucketLen*h.bucketNum)
	} else {
		h.bucketLen = blen
		h.bucketNum = bnum
		h.nodes = make([]*SHashNode, h.bucketLen*h.bucketNum)
	}
	return h
}

//
func (h *SHashTable) getnode(key string) (*SHashNode, int) {
	h.RLock()
	defer h.RUnlock()
	hashint := hashalgo.Hash33(key)
	start_index := (hashint % h.bucketNum) * h.bucketLen

	var free bool = false
	var prevNode *SHashNode
	var ret_pos int
	for i := 0; i < h.bucketLen; i++ {
		curpos := start_index + i
		pNode := h.nodes[curpos]
		if pNode == nil {
			ret_pos = curpos
			free = true
			continue
		}
		if key == pNode.key {
			//found
			return pNode, curpos
		}
		//pNode!=nil
		if free {
			continue
		}
		if prevNode == nil {
			prevNode = pNode
			ret_pos = curpos
		}
		if prevNode.vistor > pNode.vistor {
			prevNode = pNode
			ret_pos = curpos
			continue
		}
	}

	//
	return nil, ret_pos
}

func (h *SHashTable) Len() int {
	h.RLock()
	defer h.RUnlock()
	return len(h.nodes)
}

func (h *SHashTable) Get(key string) (interface{}, bool) {
	node, _ := h.getnode(key)
	if node != nil {
		return node.value, true
	}
	return nil, false
}

//LRU
func (h *SHashTable) Set(key string, value interface{}) {
	now := time.Now().Unix()
	node, position := h.getnode(key)
	if node != nil {
		//exists,update time
		h.Lock()
		node.value, node.vistor = value, now
		h.Unlock()
	} else {
		//create a new node
		h.Lock()
		h.nodes[position] = &SHashNode{
			key:    key,
			value:  value,
			vistor: now,
		}
		h.Unlock()
	}
	return
}

func (h *SHashTable) Del(key string) bool {
	node, pos := h.getnode(key)
	if node != nil {
		h.Lock()
		defer h.Unlock()
		h.nodes[pos] = nil
		return true
	}
	return false
}

/*
func main() {
	ht := (&SHashTable{}).Init(DEFAULT_BUCKET_LEN, DEFAULT_BUCKET_NUM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			ht.Set(strconv.Itoa(i), strconv.Itoa(i))
		}

	}()
	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			ht.Del(strconv.Itoa(i))
		}

	}()

	wg.Wait()
	for i := 0; i < 100000; i++ {
		_, found := ht.Get(strconv.Itoa(i))
		if found {
			fmt.Println("found")
		}
	}
	fmt.Println("len", ht.Len())
}
*/
