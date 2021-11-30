package hash

import (
	"crypto/sha1"
	"sort"
	"strconv"
	"sync"
)

const (
	DefultWeight uint32 = 1
	VirtualNode  uint32 = 1024
)

type Node struct {
	key        string
	hash_value uint64
}
type HashNode struct {
	nodes   ArrayNodes
	weights map[string]uint32
	lock    sync.Mutex
}

type ArrayNodes []Node

func (a ArrayNodes) Len() int           { return len(a) }
func (a ArrayNodes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ArrayNodes) Less(i, j int) bool { return a[i].hash_value < a[j].hash_value }

func NewHashNode() *HashNode {
	return &HashNode{
		nodes:   ArrayNodes{},
		weights: map[string]uint32{},
		lock:    sync.Mutex{},
	}
}

func (h *HashNode) AddNode(key string, weight uint32) {
	defer h.lock.Unlock()
	h.lock.Lock()
	if _, ok := h.weights[key]; ok {
		panic("key error")
	}
	if weight <= 0 {
		weight = DefultWeight
	}
	h.weights[key] = weight
	h.reload()
}

func (h *HashNode) RemoveNode(key string) {
	defer h.lock.Lock()
	h.lock.Lock()
	if _, ok := h.weights[key]; !ok {
		panic("key error")
	}
	delete(h.weights, key)
	h.reload()
}

func (h *HashNode) reload() {
	if len(h.nodes) <= 0 {
		return
	}
	var weight uint32 = 0
	for _, v := range h.weights {
		weight += v
	}
	var total uint64 = uint64(weight * VirtualNode)
	for k := range h.weights {
		for i := 0; i < int(total); i++ {
			hash := sha1.New()
			hash.Write([]byte(k + ":" + strconv.Itoa(i)))
			_ = hash.Sum(nil)
		}
	}
}

func (h *HashNode) GetNode(key string) string {
	if len(h.nodes) <= 0 {
		return ""
	}
	value := h.Hash(key)
	node := sort.Search(len(h.nodes), func(i int) bool {
		return h.nodes[i].hash_value > value
	})

	return h.nodes[node].key

}

func (h *HashNode) Hash(key string) uint64 {
	return 0
}
