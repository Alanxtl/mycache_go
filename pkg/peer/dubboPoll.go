package peer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/logger"
	"dubbo.apache.org/dubbo-go/v3/metrics"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/joho/godotenv"
)

import (
	"github.com/Alanxtl/mycache_go/pkg/message"
	"github.com/Alanxtl/mycache_go/pkg/mycache"
	"github.com/Alanxtl/mycache_go/pkg/mycache/loadbalance"
	"github.com/Alanxtl/mycache_go/pkg/mycache/loadbalance/consistenthash"
)

type DubboPoll struct {
	self         string
	basePath     string
	lock         sync.Mutex
	peers        loadbalance.Loadbalance
	dubboGetters map[string]*DubboGetter
}

func NewDubboPoll(self string) *DubboPoll {
	return &DubboPoll{
		self:     self,
		basePath: DefaultBasePath,
		lock:     sync.Mutex{},
	}
}

func (p *DubboPoll) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *DubboPoll) Get(context context.Context, req *message.Request) (*message.Response, error) {
	groupName := req.Group
	key := req.Key

	if groupName == "" || key == "" {
		return nil, fmt.Errorf("groupName or key is empty")
	}

	group := mycache.GetGroup(groupName)
	if group == nil {
		return nil, fmt.Errorf("no such group %s", groupName)
	}

	view, err := group.Get(key)
	if err != nil {
		return nil, fmt.Errorf("get group %s failed: %v", groupName, err)
	}

	return &message.Response{Value: view.ByteSlice()}, nil
}

func (p *DubboPoll) Serve(url string) {
	parts := strings.Split(url, ":")
	if len(parts) != 3 {
		return
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(parts[2])
	if err != nil {
		return
	}

	registryAddr := os.Getenv("REGISTRY_ADDR")
	PrometheusGatewayUsername := os.Getenv("PrometheusGatewayUsername")
	PrometheusGatewayPassword := os.Getenv("PrometheusGatewayPassword")
	PrometheusGatewayUrl := os.Getenv("PrometheusGatewayUrl")

	ins, err := dubbo.NewInstance(
		dubbo.WithName("mycache_server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(registryAddr),
		),
		dubbo.WithMetrics(
			metrics.WithEnabled(),                   // default false
			metrics.WithPrometheus(),                // set prometheus metric, default prometheus
			metrics.WithPrometheusExporterEnabled(), // enable prometheus exporter default false
			metrics.WithPort(9090),                  // prometheus http exporter listen at 9099,default 9090
			metrics.WithPath("/metrics"),            // prometheus http exporter url path, default /metrics
			metrics.WithMetadataEnabled(),           // enable metadata center metrics, default true
			metrics.WithRegistryEnabled(),           // enable registry metrics, default true
			metrics.WithConfigCenterEnabled(),       // enable config center metrics, default true

			metrics.WithPrometheusPushgatewayEnabled(), // enable prometheus pushgateway
			metrics.WithPrometheusGatewayUsername(PrometheusGatewayUsername),
			metrics.WithPrometheusGatewayPassword(PrometheusGatewayPassword),
			metrics.WithPrometheusGatewayUrl(PrometheusGatewayUrl), // host:port or ip:port,“http://” is added automatically,do not include the “/metrics/jobs/…” part
			metrics.WithPrometheusGatewayInterval(time.Second*10),
			metrics.WithPrometheusGatewayJob("push"), // set a metric job label, job=push to metric

			metrics.WithAggregationEnabled(), // enable rpc metrics aggregations，Most of the time there is no need to turn it on, default false
			metrics.WithAggregationTimeWindowSeconds(30),
			metrics.WithAggregationBucketNum(10), // agg bucket num
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(port),
		),
		dubbo.WithLogger(
			logger.WithLevel("warn"),
		),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()

	if err != nil {
		panic(err)
	}
	if err := message.RegisterGroupCacheHandler(srv, &DubboPoll{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		p.Log(err.Error())
	}

}

func (p *DubboPoll) Set(peers ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.dubboGetters = make(map[string]*DubboGetter, len(peers))
	for _, peer := range peers {
		p.dubboGetters[peer] = &DubboGetter{BaseURL: peer + p.basePath}
	}
}

func (p *DubboPoll) UpdatePeers(peers ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.dubboGetters = make(map[string]*DubboGetter, len(peers))
	for _, peer := range peers {
		p.dubboGetters[peer] = &DubboGetter{BaseURL: peer + p.basePath}
	}
}

func (p *DubboPoll) PickPeer(key string) (mycache.PeerGetter, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.dubboGetters[peer], true
	} else if peer == p.self {
		p.Log("Pick self, pick peer fail %s", p.self)
	}

	return nil, false
}

func (p *DubboPoll) GetSelf() string {
	return p.self
}

//var _ mycache.PeerPicker = (*DubboPoll)(nil)
