package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	pb "csrf/csrf/proto"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CSRFServiceServer struct {
	pb.UnimplementedCSRFServiceServer
	redisClient RedisClient
	secretKey   []byte
}

func NewCSRFServiceServer(redisClient RedisClient, secretKey []byte) *CSRFServiceServer {
	return &CSRFServiceServer{
		redisClient: redisClient,
		secretKey:   secretKey,
	}
}

func (s *CSRFServiceServer) GenerateToken(ctx context.Context, req *pb.Empty) (*pb.GenerateTokenResponse, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token: %v", err)
	}

	tokenString := base64.StdEncoding.EncodeToString(token)
	err = s.redisClient.Set(ctx, tokenString, "valid", time.Hour).Err()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save token to Redis: %v", err)
	}

	// Referer validation
	if req.Referer != "" {
		isTrusted := false
		for _, domain := range s.trustedDomains {
			if strings.HasPrefix(req.Referer, domain) {
				isTrusted = true
				break
			}
		}
		if !isTrusted {
			return &pb.ValidateTokenResponse{IsValid: false}, nil
		}
	}

	return &pb.GenerateTokenResponse{
		Token: tokenString,
	}, nil
}

func (s *CSRFServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	val, err := s.redisClient.Get(ctx, req.Token).Result()
	log.Println("test")
	if err == redis.Nil || val != "valid" {
		return &pb.ValidateTokenResponse{IsValid: false}, nil
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get token from Redis: %v", err)
	}

	return &pb.ValidateTokenResponse{IsValid: true}, nil
}
