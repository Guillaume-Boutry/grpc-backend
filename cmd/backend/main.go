package main

import (
	"fmt"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	service "github.com/Guillaume-Boutry/grpc-backend/pkg/services"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Config struct {
	Port int `default:"9000"`
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Starting grpc server on port %d\n", config.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
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
