package main

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/carsonoid/talk-lets-auth-with-go/pkg/simplejwt"
	"github.com/golang-jwt/jwt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE %s <public-ed-key-path>\n", os.Args[0])
		os.Exit(1)
	}

	backend, err := NewBackend()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}

	// create middleware using the given public key path
	middleware, err := simplejwt.NewMiddleware(os.Args[1])
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		// just pass our middleware function as the interceptor
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptor),
	)
	pb.RegisterGreeterServer(s, backend)

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type Backend struct {
	pb.UnimplementedGreeterServer
}

func NewBackend() (*Backend, error) {
	return &Backend{}, nil
}

// SayHello implements helloworld.GreeterServer it requires a valid token in the context
// and prints the included preferred name from the request and roles from the token
func (b *Backend) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// not having a token is now an exceptional state and we can just
	// let the context helper panic if that happens
	token := simplejwt.MustContextGetToken(ctx)

	// dig the roles from the claims
	roles := token.Claims.(jwt.MapClaims)["roles"]

	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello %s! I am the backend. You have roles %v", in.GetName(), roles),
	}, nil
}
