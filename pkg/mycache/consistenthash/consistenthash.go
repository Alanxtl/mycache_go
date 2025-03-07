package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int // 虚拟节点倍数
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hash,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add 增加节点
// keys 真实节点
func (m *Map) Add(keys ...string) {
	for _, keys := range keys {
		for i := 0; i < m.replicas; i++ {
			// 分配虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + keys)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = keys
		}
	}
	sort.Ints(m.keys)
}

// Get 传入key 输出对应的虚拟节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}

func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
