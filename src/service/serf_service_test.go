package service

import (
	"context"
	"testing"
	"time"

	pb "csrf/csrf/proto"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

type TestCase struct {
	name          string
	token         string
	value         string
	expectedValid bool
}

func TestGenerateToken(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&redis.StatusCmd{})

	secretKey := []byte("your-secret-key-must-be-32-bytes-long!")
	service := NewCSRFServiceServer(mockRedis, secretKey)

	ctx := context.Background()
	resp, err := service.GenerateToken(ctx, &pb.Empty{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Token == "" {
		t.Fatalf("expected non-empty token, got empty")
	}

	mockRedis.AssertExpectations(t)
}

func TestValidateToken(t *testing.T) {
	testCases := []TestCase{
		{
			name:          "Valid token",
			token:         "valid-token",
			value:         "valid",
			expectedValid: true,
		},
		{
			name:          "Invalid token",
			token:         "invalid-token",
			value:         "invalid",
			expectedValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := new(MockRedisClient)
			redisCmd := &redis.StringCmd{}
			redisCmd.SetVal(tc.value)

			mock.On("Get", context.Background(), tc.token).Return(redisCmd)

			secretKey := []byte("your-secret-key-must-be-32-bytes-long!")
			service := NewCSRFServiceServer(mock, secretKey)

			resp, err := service.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: tc.token})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.IsValid != tc.expectedValid {
				t.Fatalf("expected token validity to be %v, got %v", tc.expectedValid, resp.IsValid)
			}

			mock.AssertExpectations(t)
		})
	}
}
