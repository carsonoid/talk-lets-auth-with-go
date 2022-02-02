package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/carsonoid/talk-lets-auth-with-go/pkg/simplejwt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE %s <public-ed-key-path>\n", os.Args[0])
		os.Exit(1)
	}

	// create middleware using the given public key path
	middleware, err := simplejwt.NewMiddleware(os.Args[1])
	if err != nil {
		panic(err)
	}
	// END HTTP OMIT

	// create our client, add the client interceptor (middleware)
	// that way we automatically pass on the token from the context
	// to the next call
	conn, err := grpc.Dial("localhost:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	backendClient := pb.NewGreeterClient(conn)

	// note how the frontend doesn't need the key now
	frontend, err := NewFrontend(backendClient)
	if err != nil {
		panic(err)
	}

	// create "business logic" mux
	// thanks to the middleware, we can just write simple handlers
	// that can assume auth is always done
	mux := http.NewServeMux()
	// add handlers here
	mux.HandleFunc("/claims", frontend.ClaimsHandler)
	mux.HandleFunc("/hello", frontend.HelloHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n")) //nolint
	})

	// create a root mux
	root := http.NewServeMux()

	// all routes run through the middleware and call the business logic mux
	root.Handle("/", middleware.HandleHTTP(mux))

	fmt.Println("Listening on :8082")
	err = http.ListenAndServe(":8082", root)
	if err != nil {
		panic(err)
	}
}

type Frontend struct {
	backendClient pb.GreeterClient
}

func NewFrontend(backendClient pb.GreeterClient) (*Frontend, error) {
	return &Frontend{
		backendClient: backendClient,
	}, nil
}

func (f *Frontend) ClaimsHandler(w http.ResponseWriter, r *http.Request) {
	// get the token from the context to write the claims
	token := simplejwt.MustContextGetToken(r.Context())

	_, _ = w.Write([]byte(fmt.Sprint(token.Claims)))
}

func (f *Frontend) HelloHandler(w http.ResponseWriter, r *http.Request) {
	// get preferred name from the request, default to my friend
	preferredName := r.URL.Query().Get("preferredName")
	if preferredName == "" {
		preferredName = "my friend"
	}

	// make the call with the request context, the token is automatically read from the request context
	// and passed to the grpc backend though metadata
	resp, err := f.backendClient.SayHello(r.Context(), &pb.HelloRequest{Name: preferredName})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not greet: %v", err))) //nolint
		return
	}

	w.Write([]byte(fmt.Sprintf("Greeting: %s", resp.GetMessage()))) //nolint
}
