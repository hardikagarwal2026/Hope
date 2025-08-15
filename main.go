package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"hope/di"
	"hope/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	
	authv1 "hope/proto/v1/auth"
	chatv1 "hope/proto/v1/chat"
	locationv1 "hope/proto/v1/location"
	matchv1 "hope/proto/v1/match"
	reviewv1 "hope/proto/v1/review"
	ridev1 "hope/proto/v1/ride"
	userv1 "hope/proto/v1/user"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	handlers, err := di.InitApp()
	if err != nil {
		log.Fatalf("DI bootstrap failed: %v", err)
	}

	authConfig := middleware.Config{
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),
		PublicMethods: map[string]bool{
			"/proto.v1.AuthService/Login": true,
		},
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(authConfig)),
	)

	authv1.RegisterAuthServiceServer(grpcServer, handlers.AuthHandler)
	chatv1.RegisterChatServiceServer(grpcServer, handlers.ChatHandler)
	locationv1.RegisterLocationServiceServer(grpcServer, handlers.LocationHandler)
	matchv1.RegisterMatchServiceServer(grpcServer, handlers.MatchHandler)
	reviewv1.RegisterReviewServiceServer(grpcServer, handlers.ReviewHandler)
	ridev1.RegisterRideServiceServer(grpcServer, handlers.RideHandler)
	userv1.RegisterUserServiceServer(grpcServer, handlers.UserHandler)

	
	reflection.Register(grpcServer)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}
	fmt.Printf("gRPC server listening on %s\n", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server exited: %v", err)
	}
}
