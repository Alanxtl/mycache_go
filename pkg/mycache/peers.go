package mycache

import (
	pb "github.com/Alanxtl/mycache_go/pkg/message"
)

type PeerPicker interface {
	GetSelf() string
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *pb.Request) (out *pb.Response, err error)
}
