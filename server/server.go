package server

import (
	"fmt"
	"os"

	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"

	"go.uber.org/zap"

	"github.com/factorysh/jaeger-lite/conf"
	_reporter "github.com/factorysh/jaeger-lite/reporter"
	"github.com/jaegertracing/jaeger/cmd/agent/app/processors"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/cmd/agent/app/servers/thriftudp"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

type Server interface {
	Serve()
}

func NewServer(cfg *conf.Config) (Server, error) {
	//metricsFactory := metrics.NewLocalFactory(0)
	f := &Factory{}

	transport, err := thriftudp.NewTUDPServerTransport(cfg.Listen)
	if err != nil {
		return nil, err
	}
	maxPacketSize := 65000
	queueSize := 100
	server, err := servers.NewTBufferedServer(transport, queueSize, maxPacketSize, f)
	if err != nil {
		return nil, err
	}
	compactFactory := thrift.NewTCompactProtocolFactory()
	l := zap.NewExample()

	reporters := make(reporter.MultiReporter, 0)
	for name, r := range cfg.Reporters {
		f, ok := _reporter.Reporters[name]
		if !ok {
			return nil, fmt.Errorf("Unknown reporter: %v", f)
		}
		rr, err := f(r)
		if err != nil {
			return nil, err
		}
		reporters = append(reporters, rr)
	}
	handler := jaegerThrift.NewAgentProcessor(reporters)
	return processors.NewThriftProcessor(server, 1, f, compactFactory, handler, l)
}

func New() (Server, error) {
	cfgPath := os.Getenv("CONFIG")
	if cfgPath == "" {
		cfgPath = "/etc/jaeger-lite.yml"
	}
	cfg, err := conf.Read(cfgPath)
	if err != nil {
		return nil, err
	}
	return NewServer(cfg)
}
