package main

import (
	"context"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/mirrorsge/grpc-lb/balancer"
	"github.com/mirrorsge/grpc-lb/examples/proto"
	"github.com/mirrorsge/grpc-lb/registry/consul"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	consul.RegisterResolver("consul", &consulapi.Config{Address: "http://59.110.162.134:8500"}, "alpha:1.0")
	conn, err := grpc.Dial("consul:///", grpc.WithInsecure(), grpc.WithBalancerName(balancer.RoundRobin))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer conn.Close()

	client := proto.NewAlphaClient(conn)

	res, err := client.Hello(context.Background(), &proto.HelloReq{
		Name: "lijianjun",
		Time: time.Now().String(),
	})
	if err != nil {
		log.Printf("grpc err: %s", err)
		return
	}
	log.Println("result is ", res.Greeting)
}
