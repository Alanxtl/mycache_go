package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"fmt"
	message "github.com/Alanxtl/mycache_go/pkg/message"
	"github.com/dubbogo/gost/log/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"github.com/Alanxtl/mycache_go/pkg/mycache/getter"
	"github.com/Alanxtl/mycache_go/pkg/server"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *mycache.Group {
	return mycache.NewGroup("scores", 2<<10, getter.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, gee *mycache.Group) {
	peers := server.NewHttpPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *mycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
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

type GroupCacheServer struct {
}


func (srv *GroupCacheServer) Get(ctx context.Context, req *message.Request) (*message.Response, error) {
	resp := &message.Response{Value: nil}
	return resp, nil
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	registryAddr := os.Getenv("REGISTRY_ADDR")

	ins, err := dubbo.NewInstance(
		dubbo.WithName("mycache_server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(registryAddr),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err := message.RegisterGroupCacheHandler(srv, &GroupCacheServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
