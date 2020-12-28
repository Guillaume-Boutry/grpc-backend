package service

import (
	"context"
	"github.com/Guillaume-Boutry/grpc-backend/pkg/face_authenticator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type EnrollerServiceGrpcImpl struct {
	face_authenticator.UnimplementedEnrollerServer
}

//NewEnrollServiceGrpcImpl
func NewEnrollServiceGrpcImpl() *EnrollerServiceGrpcImpl {
	return &EnrollerServiceGrpcImpl{}
}

func (serviceImpl *EnrollerServiceGrpcImpl) Enroll(ctx context.Context, request *face_authenticator.EnrollRequest) (*face_authenticator.EnrollResponse, error) {
	log.Println(request.FaceRequest.Id)

	return &face_authenticator.EnrollResponse{
		Status: face_authenticator.EnrollStatus_ENROLL_STATUS_ERROR,
		Message: "ok",
	}, status.Error(codes.Canceled, "Test")
}