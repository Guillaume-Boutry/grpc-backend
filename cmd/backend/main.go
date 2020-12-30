package main

import (
	"fmt"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	service "github.com/Guillaume-Boutry/grpc-backend/pkg/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("Starting grpc server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	enrollImpl := service.NewEnrollServiceGrpcImpl()
	authenticateImpl := service.NewAuthenticateServiceGrpcImpl()

	grpcServer := grpc.NewServer()

	face_authenticator.RegisterEnrollerServer(grpcServer, enrollImpl)
	face_authenticator.RegisterAuthenticatorServer(grpcServer, authenticateImpl)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
