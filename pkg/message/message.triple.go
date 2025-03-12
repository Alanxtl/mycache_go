// Code generated by protoc-gen-triple. DO NOT EDIT.
//
// Source: message.proto
package message

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
)

// This is a compile-time assertion to ensure that this generated file and the Triple package
// are compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of Triple newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of Triple or updating the Triple
// version compiled into your binary.
const _ = triple_protocol.IsAtLeastVersion0_1_0

const (
	// GroupCacheName is the fully-qualified name of the GroupCache service.
	GroupCacheName = "message.GroupCache"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// GroupCacheGetProcedure is the fully-qualified name of the GroupCache's Get RPC.
	GroupCacheGetProcedure = "/message.GroupCache/Get"
)

var (
	_ GroupCache = (*GroupCacheImpl)(nil)
)

// GroupCache is a client for the message.GroupCache service.
type GroupCache interface {
	Get(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error)
}

// NewGroupCache constructs a client for the message.GroupCache service.
func NewGroupCache(cli *client.Client, opts ...client.ReferenceOption) (GroupCache, error) {
	conn, err := cli.DialWithInfo("message.GroupCache", &GroupCache_ClientInfo, opts...)
	if err != nil {
		return nil, err
	}
	return &GroupCacheImpl{
		conn: conn,
	}, nil
}

func SetConsumerGroupCache(srv common.RPCService) {
	dubbo.SetConsumerServiceWithInfo(srv, &GroupCache_ClientInfo)
}

// GroupCacheImpl implements GroupCache.
type GroupCacheImpl struct {
	conn *client.Connection
}

func (c *GroupCacheImpl) Get(ctx context.Context, req *Request, opts ...client.CallOption) (*Response, error) {
	resp := new(Response)
	if err := c.conn.CallUnary(ctx, []interface{}{req}, resp, "Get", opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

var GroupCache_ClientInfo = client.ClientInfo{
	InterfaceName: "message.GroupCache",
	MethodNames:   []string{"Get"},
	ConnectionInjectFunc: func(dubboCliRaw interface{}, conn *client.Connection) {
		dubboCli := dubboCliRaw.(*GroupCacheImpl)
		dubboCli.conn = conn
	},
}

// GroupCacheHandler is an implementation of the message.GroupCache service.
type GroupCacheHandler interface {
	Get(context.Context, *Request) (*Response, error)
}

func RegisterGroupCacheHandler(srv *server.Server, hdlr GroupCacheHandler, opts ...server.ServiceOption) error {
	return srv.Register(hdlr, &GroupCache_ServiceInfo, opts...)
}

func SetProviderGroupCache(srv common.RPCService) {
	dubbo.SetProviderServiceWithInfo(srv, &GroupCache_ServiceInfo)
}

var GroupCache_ServiceInfo = server.ServiceInfo{
	InterfaceName: "message.GroupCache",
	ServiceType:   (*GroupCacheHandler)(nil),
	Methods: []server.MethodInfo{
		{
			Name: "Get",
			Type: constant.CallUnary,
			ReqInitFunc: func() interface{} {
				return new(Request)
			},
			MethodFunc: func(ctx context.Context, args []interface{}, handler interface{}) (interface{}, error) {
				req := args[0].(*Request)
				res, err := handler.(GroupCacheHandler).Get(ctx, req)
				if err != nil {
					return nil, err
				}
				return triple_protocol.NewResponse(res), nil
			},
		},
	},
}
