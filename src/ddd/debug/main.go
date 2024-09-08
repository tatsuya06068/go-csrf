package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "path/to/your/proto/generated/csrf"
)

func main() {
	// gRPCサーバーに接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// CsrfServiceのクライアントを作成
	client := pb.NewCsrfServiceClient(conn)

	// コンテキストとタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// トークンの発行をリクエスト
	sessionID := "12345"
	refer := "https://example.com"

	generateResp, err := client.GenerateToken(ctx, &pb.GenerateTokenRequest{
		SessionId: sessionID,
		Refer:     refer,
	})
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	log.Printf("Generated Token: %s", generateResp.CsrfToken)

	// トークンの検証をリクエスト
	validateResp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		CsrfToken: generateResp.CsrfToken,
		SessionId: sessionID,
		Refer:     refer,
	})
	if err != nil {
		log.Fatalf("Failed to validate token: %v", err)
	}
	log.Printf("Token is valid: %t", validateResp.IsValid)
}