package server

import (
	"fmt"
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"log"
	"net/http"
	"strings"
)

const DefaultBasePath = "/_mycache/"

type HttpPoll struct {
	self     string
	basePath string
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

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(view.ByteSlice())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
