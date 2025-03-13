package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

import (
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"github.com/Alanxtl/mycache_go/pkg/mycache/getter"
	"github.com/Alanxtl/mycache_go/pkg/peer"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *mycache.Group {
	return mycache.NewGroup("scores", 2<<10, getter.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *mycache.Group) {
	peers := peer.NewDubboPoll(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	peers.Serve(addr)
	log.Println("geecache is running at", addr)
	//log.Fatal()
}

func startAPIServer(apiAddr string, gee *mycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			log.Printf("[API server] search key %s", key)
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer("http://localhost:"+strconv.Itoa(port), addrs, gee)
}
