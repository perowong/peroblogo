package consul

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type consulBuilder struct {
}

type consulResolver struct {
	cc                   resolver.ClientConn
	ctx                  context.Context
	cancel               context.CancelFunc
	name                 string
	tag                  string
	disableServiceConfig bool
	lastIndex            uint64
}

const (
	defaultPort    = "8500"
	resolverPrefix = "consul resolver: "
)

var (
	errInvalidName    = errors.New(resolverPrefix + "invalid service name")
	errInvalidAddress = errors.New(resolverPrefix + "invalid address")
	addressRegex      = regexp.MustCompile(`^([0-9]{1,3}(\.[0-9]{1,3}){3})(:([0-9]{2,5}))?$`)
	endpointRegex     = regexp.MustCompile(`^([a-z][a-z0-9]*(-[a-z0-9]+)*)(#(.+))?$`)
)

func init() {
	resolver.Register(NewBuilder())
}

// NewBuilder return a consul resolver builder
func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	host, port, name, tag, err := parseTarget(target.Authority, target.Endpoint)
	if err != nil {
		return nil, err
	}

	cr := &consulResolver{
		cc:                   cc,
		name:                 name,
		tag:                  tag,
		disableServiceConfig: opts.DisableServiceConfig,
	}
	cr.ctx, cr.cancel = context.WithCancel(context.Background())

	go cr.watcher(host + ":" + port)
	return cr, nil
}

func (cr *consulResolver) watcher(address string) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		cr.cc.ReportError(fmt.Errorf(resolverPrefix+"%v", err))
		return
	}

	for {
		qo := &api.QueryOptions{
			WaitIndex: cr.lastIndex,
		}
		qo = qo.WithContext(cr.ctx)
		entries, meta, err := client.Health().Service(cr.name, cr.tag, true, qo)
		select {
		default:
		case <-cr.ctx.Done():
			return
		}
		if err != nil {
			cr.cc.ReportError(fmt.Errorf(resolverPrefix+"%v", err))
			continue
		}
		cr.lastIndex = meta.LastIndex

		addresses := make([]resolver.Address, 0, len(entries))
		for _, ent := range entries {
			addr := ent.Service.Address + ":" + strconv.FormatInt(int64(ent.Service.Port), 10)
			addresses = append(addresses, resolver.Address{Addr: addr})
		}
		cr.cc.UpdateState(resolver.State{
			Addresses: addresses,
		})
	}
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
}

func (cr *consulResolver) Close() {
	cr.cancel()
}

func parseTarget(address, endpoint string) (host, port, name, tag string, err error) {
	if !addressRegex.MatchString(address) {
		err = errInvalidAddress
		return
	}
	if !endpointRegex.MatchString(endpoint) {
		err = errInvalidName
		return
	}

	groups := addressRegex.FindStringSubmatch(address)
	host = groups[1]
	port = groups[4]
	groups = endpointRegex.FindStringSubmatch(endpoint)
	name = groups[1]
	tag = groups[4]
	if port == "" {
		port = defaultPort
	}
	return
}
