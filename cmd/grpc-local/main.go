package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// make the call without any custom context
	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "local"})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Message)
}
