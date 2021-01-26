package main

import (
	"context"
	"github.com/mirrorsge/grpc-lb/examples/proto"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	res, err := proto.AlphaAdapter.Hello(ctx, &proto.HelloReq{
		Name: "lijianjun",
		Time: time.Now().String(),
	})
	if err != nil {
		log.Printf("grpc err: %s", err)
		return
	}
	log.Println("result is ", res.Greeting)
}
