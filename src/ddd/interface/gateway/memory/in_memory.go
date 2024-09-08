package memory

import (
	"context"
	"fmt"
	"time"

	pb "csrf/csrf/proto"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"csrf/service"

	csrf "csrf/ddd/domain"
)

type CSRFGateway struct {
	redisClient service.RedisClient
}

func NewCSRFGateway(redisClient service.RedisClient) *CSRFGateway {
	return &CSRFGateway{
		redisClient: redisClient,
	}
}

func (cg *CSRFGateway) Save(ctx context.Context, csrfToken csrf.CsrfToken) error {

	key := fmt.Sprintf("csrf_token_%s", csrfToken.SessionID)

	err := cg.redisClient.Set(ctx, csrfToken.Token, key, time.Hour).Err()
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to save token to Redis: %v", err)
	}

	return nil
}

func (cg *CSRFGateway) Find(ctx context.Context, req *pb.ValidateTokenRequest) (string, error) {
	val, err := cg.redisClient.Get(ctx, req.Token).Result()

	if err == redis.Nil {
		return "", status.Error(codes.Internal, "Failed to get token from Redis: nil")
	} else if err != nil {
		return "", status.Errorf(codes.Internal, "Failed to get token from Redis: %v", err)
	}

	return val, nil
}
