package mycache

import (
	"fmt"
	"github.com/Alanxtl/mycache_go/pkg/cache"
	"github.com/Alanxtl/mycache_go/pkg/cache/lru"
	"github.com/Alanxtl/mycache_go/pkg/mycache/getter"
	"github.com/Alanxtl/mycache_go/pkg/tools"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    getter.Getter
	mainCache lru.Mutex[string, cache.ByteView]
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int, getter getter.Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: lru.Mutex[string, cache.ByteView]{
			Cap: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}

func (g *Group) Get(key string) (cache.ByteView, error) {
	if key == "" {
		return cache.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		// 从 mainCache 中查找缓存，如果存在则返回缓存值
		log.Println("[mycache] hit")
		return v, nil
	}

	// 缓存不存在，则调用 load 方法
	return g.load(key)
}

func (g *Group) load(key string) (cache.ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (cache.ByteView, error) {
	// 通过用户回调函数获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return cache.ByteView{}, err
	}

	value := cache.ByteView{Bytes: tools.CloneBytes(bytes)}

	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value cache.ByteView) {
	// 将源数据添加到缓存
	g.mainCache.Add(key, value)
}
