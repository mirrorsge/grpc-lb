package consul

import (
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"sync"
)

type consulResolver struct {
	scheme      string
	consulConf  *consulapi.Config
	ServiceName string
	watcher     *ConsulWatcher
	cc          resolver.ClientConn
	wg          sync.WaitGroup
}

func (r *consulResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	r.watcher = newConsulWatcher(r.ServiceName, r.consulConf)
	r.start()
	return r, nil
}

func (r *consulResolver) Scheme() string {
	return r.scheme
}

func (r *consulResolver) start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		out := r.watcher.Watch()
		for addr := range out {
			r.cc.UpdateState(resolver.State{Addresses: addr})
		}
	}()
}

func (r *consulResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *consulResolver) Close() {
	r.watcher.Close()
	r.wg.Wait()
}

func RegisterResolver(scheme string, consulConf *consulapi.Config, srvName string) {
	resolver.Register(&consulResolver{
		scheme:      scheme,
		consulConf:  consulConf,
		ServiceName: srvName,
	})
}
