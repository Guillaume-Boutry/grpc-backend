package service

import (
	"context"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type AuthenticateServiceGrpcImpl struct {
	face_authenticator.UnimplementedAuthenticatorServer
	client cloudevents.Client
	Target string `envconfig:"K_SINK_AUTHENTICATE"`
}

//NewEnrollServiceGrpcImpl
func NewAuthenticateServiceGrpcImpl() *AuthenticateServiceGrpcImpl {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	authService :=  &AuthenticateServiceGrpcImpl{
		client: client,
	}
	if err := envconfig.Process("", authService); err != nil {
		log.Fatal(err.Error())
	}
	return authService
}


func (serviceImpl *AuthenticateServiceGrpcImpl) Authenticate(ctx context.Context, request *face_authenticator.AuthenticateRequest) (*face_authenticator.AuthenticateResponse, error) {
	log.Println(request.FaceRequest.Id)

	return &face_authenticator.AuthenticateResponse{
		Status:  face_authenticator.AuthenticateStatus_AUTHENTICATE_STATUS_ERROR,
		Message: "ok",
	}, nil
}
