package peer

import (
	"context"
	"log"
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/joho/godotenv"
)

import (
	"github.com/Alanxtl/mycache_go/pkg/message"
)

type DubboGetter struct {
	BaseURL string
}

func (h *DubboGetter) Get(in *message.Request) (*message.Response, error) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	registryAddr := os.Getenv("REGISTRY_ADDR")


	ins, err := dubbo.NewInstance(
		dubbo.WithName("mycache_client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(registryAddr),
		),
		dubbo.WithLogger(
			logger.WithLevel("warn"),
		),
	)
	if err != nil {
		panic(err)
	}

	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := message.NewGroupCache(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Get(context.Background(),
		&message.Request{
			Group: in.Group,
			Key:   in.Key})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("[getter] get response from peer %s: %s", h.BaseURL, resp)

	return resp, nil
}

//var _ mycache.PeerGetter = (*DubboGetter)(nil)
