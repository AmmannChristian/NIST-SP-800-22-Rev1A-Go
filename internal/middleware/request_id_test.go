package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestUnaryRequestIDInterceptor(t *testing.T) {
	interceptor := UnaryRequestIDInterceptor()

	// Mock handler that captures the context
	var capturedCtx context.Context
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		capturedCtx = ctx
		return "response", nil
	}

	// Create mock info
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	// Call the interceptor
	_, err := interceptor(context.Background(), "request", info, handler)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}

	// Verify request ID was added to context
	requestID := GetRequestID(capturedCtx)
	if requestID == "" {
		t.Error("request ID not found in context")
	}

	// Verify it's a valid UUID format (simple check)
	if len(requestID) != 36 {
		t.Errorf("request ID has invalid format: %s", requestID)
	}
}

func TestGetRequestID_NoID(t *testing.T) {
	ctx := context.Background()
	requestID := GetRequestID(ctx)
	if requestID != "" {
		t.Errorf("expected empty string, got %s", requestID)
	}
}

func TestGetRequestID_WithID(t *testing.T) {
	expectedID := "test-request-id-123"
	ctx := context.WithValue(context.Background(), RequestIDKey, expectedID)

	requestID := GetRequestID(ctx)
	if requestID != expectedID {
		t.Errorf("expected %s, got %s", expectedID, requestID)
	}
}
