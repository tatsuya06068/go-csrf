package main

import (
	"log"
	"net"

	pb "path/to/your/generated/proto" // 生成されたprotoファイルのパスに置き換えてください
	"your_project/csrf"               // csrf_service.go のパッケージ名に置き換えてください

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	secretKey := []byte("your-secret-key-must-be-32-bytes-long!")

	server := grpc.NewServer()
	pb.RegisterCSRFServiceServer(server, csrf.NewCSRFServiceServer(redisClient, secretKey))

	log.Println("gRPC server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
