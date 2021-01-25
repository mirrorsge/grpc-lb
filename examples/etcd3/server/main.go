package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/mirrorsge/grpc-lb/examples/proto"
	"github.com/mirrorsge/grpc-lb/registry"
	etcd "github.com/mirrorsge/grpc-lb/registry/etcd3"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var nodeID = flag.String("node", "node1", "node ID")
var port = flag.Int64("port", 8080, "listening port")

type RpcServer struct {
	proto.UnimplementedAlphaServer
	addr string
	s    *grpc.Server
}

func NewRpcServer(addr string) *RpcServer {
	return &RpcServer{
		addr: addr,
		s:    grpc.NewServer(),
	}
}

func (s *RpcServer) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("rpc listening on:%s", s.addr)
	proto.RegisterAlphaServer(s.s, s)
	_ = s.s.Serve(listener)
}

func (s *RpcServer) Stop() {
	s.s.GracefulStop()
}

func (s *RpcServer) Hello(ctx context.Context, req *proto.HelloReq) (*proto.HelloRes, error) {
	text := req.Name + req.Time
	return &proto.HelloRes{Greeting: text}, nil
}

func StartService() {
	etcdConfg := clientv3.Config{
		Endpoints: []string{"http://59.110.162.134:2379"},
	}
	service := &registry.ServiceInfo{
		InstanceId: *nodeID,
		Name:       "alpha",
		Version:    "1.0",
		Address:    fmt.Sprintf("127.0.0.1:%d", *port),
	}
	registerIns, err := etcd.NewRegisterIns(
		&etcd.Config{
			EtcdConfig:  etcdConfg,
			RegistryDir: "/backend/services",
			Ttl:         10 * time.Second,
		})
	if err != nil {
		log.Panic(err)
		return
	}
	server := NewRpcServer(fmt.Sprintf("localhost:%d", *port))

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		server.Run()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		_ = registerIns.Register(service)
		wg.Done()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	registerIns.UnRegister(service)
	server.Stop()
	wg.Wait()
}

func main() {
	flag.Parse()
	StartService()
}
