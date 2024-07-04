package main

import (
    "log"
    "net"

    "github.com/djg3577/STRIVEBackend/pkg/user"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    user.RegisterUserServiceServer(grpcServer, user.NewUserService())

    log.Printf("Starting User Service on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
