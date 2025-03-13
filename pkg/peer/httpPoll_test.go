package peer

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

import (
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"github.com/Alanxtl/mycache_go/pkg/mycache/getter"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestServer(t *testing.T) {
	mycache.NewGroup("scores", 2<<10, getter.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHttpPool(addr)
	log.Println("mycache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
