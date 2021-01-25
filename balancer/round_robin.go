package balancer

import (
	"github.com/mirrorsge/grpc-lb/common"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"math/rand"
	"sync"
)

const RoundRobin = "round_robin_x"

type roundRobinPickerBuilder struct{}

func newRoundRobinBuilder() balancer.Builder {
	return base.NewBalancerBuilder(RoundRobin, &roundRobinPickerBuilder{}, base.Config{HealthCheck: true})
}

func (b roundRobinPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	grpclog.Infof("roundrobinPicker: newPicker called with buildInfo: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, 0)
	for subConn, subConnInfo := range info.ReadySCs {
		weight := common.GetWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &roundRobinPicker{
		subConns: scs,
		next:     rand.Intn(len(scs)),
	}

}

type roundRobinPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
	next     int
}

func (p *roundRobinPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	ret := balancer.PickResult{}
	p.mu.Lock()
	ret.SubConn = p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return ret, nil
}

func init() {
	balancer.Register(newRoundRobinBuilder())
}
