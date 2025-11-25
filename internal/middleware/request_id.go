package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// UnaryRequestIDInterceptor adds a unique request ID to each gRPC request
func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Generate unique request ID
		requestID := uuid.New().String()

		// Add to context for internal use
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		// Add to outgoing metadata for clients to receive
		md := metadata.Pairs("x-request-id", requestID)
		if err := grpc.SetHeader(ctx, md); err != nil {
			// Continue even if header setting fails
			// This is not critical for request processing
		}

		return handler(ctx, req)
	}
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
