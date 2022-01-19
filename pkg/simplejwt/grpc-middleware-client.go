package simplejwt

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (m *Middleware) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Get the token from the context
	token, err := ContextGetToken(ctx)
	if err != nil {
		return fmt.Errorf("token not set in context: %w", err)
	}

	// add the auth token to the outgoing grpc context using
	// the generic grpc metadata tools
	ctx = metadata.NewOutgoingContext(ctx,
		metadata.New(
			map[string]string{
				"jwt": token.Raw,
			},
		),
	)

	fmt.Println("* gRPC CLIENT middleware set token")

	// call the invoker with everythign else untouched
	return invoker(ctx, method, req, reply, cc, opts...)
}
