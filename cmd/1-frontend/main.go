package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// GRPC OMIT
	conn, err := grpc.Dial("localhost:8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	backendClient := pb.NewGreeterClient(conn)

	frontend, err := NewFrontend(validator, backendClient)
	if err != nil {
		panic(err)
	}
	// GRPC END OMIT

	mux := http.NewServeMux()
	mux.HandleFunc("/", frontend.RootHandler)
	mux.HandleFunc("/claims", frontend.ClaimsHandler)
	mux.HandleFunc("/hello", frontend.HelloHandler)

	fmt.Println("Listening on :8082")
	err = http.ListenAndServe(":8082", mux)
	if err != nil {
		panic(err)
	}
}

type Frontend struct {
	validator     *simplejwt.Validator
	backendClient pb.GreeterClient
}

func NewFrontend(validator *simplejwt.Validator, backendClient pb.GreeterClient) (*Frontend, error) {
	return &Frontend{
		validator:     validator,
		backendClient: backendClient,
	}, nil
}

func (f *Frontend) ClaimsHandler(w http.ResponseWriter, r *http.Request) {
	// get the token so we can use it to print claims
	token, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) //nolint
		return
	}

	_, _ = w.Write([]byte(fmt.Sprint(token.Claims)))
}

func (f *Frontend) HelloHandler(w http.ResponseWriter, r *http.Request) {
	// get preferred name from the request, default to my friend
	preferredName := r.URL.Query().Get("preferredName")
	if preferredName == "" {
		preferredName = "my friend"
	}

	// get the token to pass it down, even though we don't use
	// it here, we do require it
	// START TOKEN OMIT
	token, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) //nolint
		return
	}

	// add the auth token to the outgoing grpc context using
	// the generic grpc metadata tools
	ctx := metadata.NewOutgoingContext(
		r.Context(),
		metadata.New(
			map[string]string{
				"jwt": token.Raw,
			},
		),
	)
	// SETTER END OMIT

	// make the call with the new context
	resp, err := f.backendClient.SayHello(ctx, &pb.HelloRequest{
		Name: preferredName,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not greet: %v", err))) //nolint
		return
	}

	w.Write([]byte(fmt.Sprintf("Greeting: %s", resp.GetMessage()))) //nolint
}

func (f *Frontend) RootHandler(w http.ResponseWriter, r *http.Request) {
	// get the token just to do auth, ignore the actual value
	_, err := f.getHeaderToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("auth error:" + err.Error())) //nolint
		return
	}

	w.Write([]byte("ok\n")) //nolint
}

// getHeaderToken checks for a valid JWT token in the Authorization header
func (f *Frontend) getHeaderToken(h http.Header) (*jwt.Token, error) {
	auth := strings.Split(h.Get("Authorization"), " ")
	if len(auth) < 2 || auth[0] != "Bearer" {
		return nil, errors.New("invalid Authorization header")
	}
	tokenString := auth[1]

	// Parse the token from the header
	token, err := f.validator.GetToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}
