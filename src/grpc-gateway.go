package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"your-module-path/internal/gateway"
	"your-module-path/internal/module1"
	"your-module-path/internal/module2"
	pb "your-module-path/pb"

	"google.golang.org/grpc"
)

func main() {
	grpcAddress := ":50051"
	httpAddress := ":8080"

	// gRPC サーバーの起動
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterModule1ServiceServer(grpcServer, &module1.Module1Service{})
	pb.RegisterModule2ServiceServer(grpcServer, &module2.Module2Service{})

	go func() {
		log.Println("Starting gRPC server on", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// gRPC-Gateway サーバーの起動
	go func() {
		log.Println("Starting HTTP/1.1 REST server on", httpAddress)
		if err := gateway.RunGatewayServer(grpcAddress, httpAddress); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// シグナルを待機してサーバーをシャットダウン
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
	os.Exit(0)
}
