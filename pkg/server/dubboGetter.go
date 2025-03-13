package server

import (
	"context"
	"log"
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"

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

	log.Println(registryAddr)

	ins, err := dubbo.NewInstance(
		dubbo.WithName("mycache_client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(registryAddr),
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

	log.Println(2)
	svc, err := message.NewGroupCache(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Get(context.Background(),
		&message.Request{
			Group: in.Group,
			Key:   in.Key})
	log.Println(3)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	log.Printf("[getter] get response: %s", resp)

	return resp, nil
}

//var _ mycache.PeerGetter = (*DubboGetter)(nil)
