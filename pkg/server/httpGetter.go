package server

import (
	"fmt"
	message "github.com/Alanxtl/mycache_go/pkg/message"
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"net/url"
)

type HttpGetter struct {
	BaseURL string
}

func (h *HttpGetter) Get(in *message.Request) (*message.Response, error) {

	u := fmt.Sprintf(
		"%v%v/%v",
		h.BaseURL,
		url.QueryEscape(in.Group),
		url.QueryEscape(in.Key),
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error: %s", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body err: %v", err)
	}

	out := &message.Response{}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return nil, fmt.Errorf("decode response body err: %v", err)
	}

	return out, nil
}

var _ mycache.PeerGetter = (*HttpGetter)(nil)
