package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"

	"github.com/carsonoid/talk-lets-auth-with-go/pkg/simplejwt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE %s <public-ed-key-path>\n", os.Args[0])
		os.Exit(1)
	}

	validator, err := simplejwt.NewValidator(os.Args[1])
	if err != nil {
		panic(err)
	}

	// PG2 OMIT
	backend, err := NewBackend(validator)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, backend)

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

type Backend struct {
	pb.UnimplementedGreeterServer
	validator *simplejwt.Validator
}

func NewBackend(validator *simplejwt.Validator) (*Backend, error) {
	return &Backend{
		validator: validator,
	}, nil
}

// SayHello implements helloworld.GreeterServer it requires a valid
// token in the context and prints the included preferred name from the request
// and roles from the token
func (b *Backend) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	token, err := b.tokenFromContextMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get token: %w", err)
	}

	// dig the roles from the claims
	roles := token.Claims.(jwt.MapClaims)["roles"]

	return &pb.HelloReply{
		Message: fmt.Sprintf(
			"Hello %s! I am the backend. You have roles %v",
			in.GetName(), roles),
	}, nil
}

func (b *Backend) tokenFromContextMetadata(ctx context.Context) (*jwt.Token, error) {
	// rip the token from the metadata via the context
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata found in context")
	}
	tokens := headers.Get("jwt")
	if len(tokens) < 1 {
		return nil, errors.New("no token found in metadata")
	}
	tokenString := tokens[0]

	token, err := b.validator.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
