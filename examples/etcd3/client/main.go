package main

import (
	"context"
	"github.com/mirrorsge/grpc-lb/balancer"
	"github.com/mirrorsge/grpc-lb/examples/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("etcd3:///", grpc.WithInsecure(), grpc.WithBalancerName(balancer.RoundRobin))
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
