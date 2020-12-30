package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/protobuf/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if facereq := request.FaceRequest; facereq != nil {
		if face := facereq.Face; face != nil {
			r := cloudevents.NewEvent(cloudevents.VersionV1)
			r.SetType("authenticate")
			r.SetSource("grpc-backend")
			binary, err := proto.Marshal(request)
			if  err != nil {
				return nil, status.Error(codes.Internal, "Fail to serialize request")
			}
			req := &Message{Payload: binary}
			if err := r.SetData("application/json", req); err != nil {
				return nil, status.Error(codes.Internal, "Failed to send request to enroller")
			}
			newCtx := cloudevents.ContextWithTarget(ctx, serviceImpl.Target)

			result, err := serviceImpl.client.Request(newCtx, r)
			// Reçu depuis Enroller
			if cloudevents.IsNACK(err) {
				return nil, status.Error(codes.Internal, fmt.Sprintf("Authenticator failed to perform task: %#v, %#v", err, errors.Unwrap(err)))
			}

			response := &face_authenticator.AuthenticateResponse{}
			if err := proto.Unmarshal(result.Data(), response); err != nil {
				return nil, status.Error(codes.Internal, "Failed to parse authenticator response")
			}
			// Envoie au client
			return response, nil
		}
	}
	return nil, status.Error(codes.Canceled, "FaceRequest cannot be empty")

}
