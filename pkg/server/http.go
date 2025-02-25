package server

import (
	"fmt"
	"github.com/Alanxtl/mycache_go/pkg/client"
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"github.com/Alanxtl/mycache_go/pkg/mycache/consistenthash"
	pb "github.com/Alanxtl/mycache_go/pkg/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	DefaultBasePath = "/_mycache/"
	defaultReplicas = 50
)

type HttpPoll struct {
	self        string
	basePath    string
	lock        sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*client.HttpGetter
}

func NewHttpPool(self string) *HttpPoll {
	return &HttpPoll{
		self:     self,
		basePath: DefaultBasePath,
	}
}

func (p *HttpPoll) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HttpPoll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := mycache.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *HttpPoll) Set(peers ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*client.HttpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &client.HttpGetter{BaseURL: peer + p.basePath}
	}
}

func (p *HttpPoll) PickPeer(key string) (mycache.PeerGetter, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _ mycache.PeerPicker = (*HttpPoll)(nil)
