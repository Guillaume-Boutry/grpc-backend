package service

import (
	"context"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/protobuf/proto"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type EnrollerServiceGrpcImpl struct {
	face_authenticator.UnimplementedEnrollerServer
	client cloudevents.Client
	Target string `envconfig:"K_SINK_ENROLL"`
}

//NewEnrollServiceGrpcImpl
func NewEnrollServiceGrpcImpl() *EnrollerServiceGrpcImpl {
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	enrollService :=  &EnrollerServiceGrpcImpl{
		client: client,
	}
	if err := envconfig.Process("", enrollService); err != nil {
		log.Fatal(err.Error())
	}
	return enrollService
}

type Message struct {
	Payload []byte `json:"payload"`
}

func (serviceImpl *EnrollerServiceGrpcImpl) Enroll(ctx context.Context, request *face_authenticator.EnrollRequest) (*face_authenticator.EnrollResponse, error) {
	log.Println(request.FaceRequest.Id)

	if facereq := request.FaceRequest; facereq != nil {
		if face := facereq.Face; face != nil {
			r := cloudevents.NewEvent(cloudevents.VersionV1)
			r.SetType("enroll")
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

			response, res := serviceImpl.client.Request(newCtx, r)
			if cloudevents.IsUndelivered(res) {
				log.Printf("Failed to request: %v", res)
				return nil, res
			} else if response != nil {
				log.Printf("Got Event Response Context: %+v\n", response.Context)
			} else {
				log.Printf("Event sent at %s", time.Now())
				return nil, status.Error(codes.Internal, "Enroller failed to perform task")
			}

			msg := &Message{}
			if err := response.DataAs(msg); err != nil {
				return nil, status.Error(codes.Internal, "Failed to parse json response")
			}
			if msg.Payload == nil {
				return nil, status.Error(codes.DataLoss, "Enroll response empty")
			}

			responseObject := &face_authenticator.EnrollResponse{}
			if err := proto.Unmarshal(msg.Payload, responseObject); err != nil {
				return nil, status.Error(codes.Internal, "Failed to parse enroller response")
			}
			// Envoie au client
			return responseObject, nil
		}
	}
	return nil, status.Error(codes.Canceled, "FaceRequest cannot be empty")
}